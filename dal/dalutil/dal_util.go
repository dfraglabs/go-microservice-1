package dalutil

import (
	"fmt"
	"strings"
	"sync"
	"time"

	newrelic "github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("pkg", "dalutil")
}

type CacheStats struct {
	hitLock sync.Locker
	hit     int

	missLock sync.Locker
	miss     int
}

func NewCacheStats() *CacheStats {
	return &CacheStats{
		hitLock:  &sync.Mutex{},
		missLock: &sync.Mutex{},
	}
}

func (c *CacheStats) RecordHit() {
	c.hitLock.Lock()
	defer c.hitLock.Unlock()

	c.hit++
}

func (c *CacheStats) RecordMiss() {
	c.missLock.Lock()
	defer c.missLock.Unlock()

	c.miss++
}

func (c *CacheStats) GetStats() map[string]int {
	stats := map[string]int{}
	c.hitLock.Lock()
	stats["hit"] = c.hit
	c.hitLock.Unlock()

	c.missLock.Lock()
	stats["miss"] = c.miss
	c.missLock.Unlock()

	return stats
}

/*****************
 Smart Collection
*****************/

type SmartCollection struct {
	coll *mgo.Collection
	mu   RWLocker
	last time.Time
	freq time.Duration
}

func NewSmartCollection(c *mgo.Collection, freq time.Duration) *SmartCollection {
	return &SmartCollection{
		coll: c,
		mu:   &sync.RWMutex{},
		last: time.Now(),
		freq: freq,
	}
}

func (s *SmartCollection) Collection() *mgo.Collection {
	s.mu.RLock()
	elapsed := time.Since(s.last)
	s.mu.RUnlock()

	if elapsed > s.freq {
		s.mu.Lock()
		s.last = time.Now()
		s.mu.Unlock()

		// this is safe to do without a lock because it implements its own lock
		s.coll.Database.Session.Refresh()
	}

	return s.coll
}

func (s *SmartCollection) EnsureIndexes(idxs []*mgo.Index) error {
	for _, idx := range idxs {
		log.Debugf("Ensuring index: %s", idx.Name)
		if err := s.UpsertIndex(idx); err != nil {
			return fmt.Errorf("Could not ensure indexes on DB: %v", err)
		}
	}

	return nil
}

// Ensure new index. If index already exists with same options, remove it and add new one.
func (s *SmartCollection) UpsertIndex(idx *mgo.Index) error {
	if err := s.coll.EnsureIndex(*idx); err != nil {
		if strings.Contains(err.Error(), "already exists with different options") ||
			strings.Contains(err.Error(), "Trying to create an index with same name") {
			log.Warnf("index already exists with name '%s'. replacing...", idx.Name)

			//drop that one
			if err := s.coll.DropIndexName(idx.Name); err != nil {
				return fmt.Errorf("failed to remove old index: %v", err)
			}

			if err := s.coll.EnsureIndex(*idx); err != nil {
				return fmt.Errorf("failed to add new index: %v", err)
			}

			return nil
		}

		return fmt.Errorf("failed to ensure index: %v", err)
	}

	return nil
}

func (s *SmartCollection) StartMongoDatastoreSegment(txn newrelic.Transaction, op string, query map[string]interface{}) *newrelic.DatastoreSegment {
	return &newrelic.DatastoreSegment{
		StartTime:       newrelic.StartSegmentNow(txn),
		Product:         newrelic.DatastoreMongoDB,
		DatabaseName:    s.coll.Database.Name,
		Collection:      s.coll.Name,
		Operation:       op,
		QueryParameters: query,
	}
}

//go:generate counterfeiter -o ../../fakes/rwlocker/rwlocker.go . RWLocker

type RWLocker interface {
	RLock()
	RUnlock()
	Lock()
	Unlock()
}
