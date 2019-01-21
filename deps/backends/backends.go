package backends

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"github.com/InVisionApp/go-health"

	"github.com/dfraglabs/go-microservice-1/config"
	"github.com/dfraglabs/go-microservice-1/dal/foo/client"
)

var (
	log *logrus.Entry
)

func init() {
	log = logrus.WithField("pkg", "backends")
}

type Backends struct {
	//Direct connection backends
	MongoDB *mgo.Database

	//API clients
	FooClient client.IClient

	//backends dependency health
	Statuses []*health.Config

	//Use this to determine if the connect method has been
	// run successfully and it is safe to use the underlying backends.
	// Not the best way to do this but good enough for now.
	connected bool
}

func NewBackends(cfg *config.Config) (*Backends, error) {
	b := &Backends{
		Statuses:  []*health.Config{},
		connected: false, //ensure
	}

	//Connect to DBs
	mongoDB, err := connectMongoDB(&MongoConfig{
		Hosts:      cfg.MongoDBHosts,
		Name:       cfg.MongoDBName,
		ReplicaSet: cfg.MongoDBReplicaSet,
		Source:     cfg.MongoDBSource,
		User:       cfg.MongoDBUser,
		Password:   cfg.MongoDBPassword,
		Timeout:    time.Duration(cfg.MongoDBConnTimeoutSec) * time.Second,
		UseSSL:     cfg.MongoDBConnUseSSL,
	})
	if err != nil {
		return nil, err
	}

	b.MongoDB = mongoDB

	// Foo client setup
	fc := client.NewFooClient(cfg.FooAPIHost, cfg.ServiceName)
	b.FooClient = fc
	b.Statuses = append(b.Statuses, &health.Config{
		Name:     "foo-client",
		Checker:  fc,
		Interval: time.Duration(cfg.HealthFreqSec) * time.Second,
	})

	b.connected = true

	return b, nil
}

type MongoConfig struct {
	Hosts      []string
	Name       string
	ReplicaSet string
	Source     string
	User       string
	Password   string
	Timeout    time.Duration
	UseSSL     bool
}

func connectMongoDB(cfg *MongoConfig) (*mgo.Database, error) {
	log.Infof("Connecting to DB: %q hosts: %v with timeout %d sec", cfg.Name, cfg.Hosts, cfg.Timeout)
	log.Debugf("DB name: '%s'; replica set: '%s'; auth source: '%s'; user: '%s'; pass len: %d; use SSL: %v",
		cfg.Name, cfg.ReplicaSet, cfg.Source, cfg.User, len(cfg.Password), cfg.UseSSL)

	dialInfo := &mgo.DialInfo{
		Addrs:          cfg.Hosts,
		Database:       cfg.Name,
		ReplicaSetName: cfg.ReplicaSet,
		Source:         cfg.Source,
		Username:       cfg.User,
		Password:       cfg.Password,
		Timeout:        cfg.Timeout,
	}

	if cfg.UseSSL {
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), &tls.Config{})
			if conn != nil {
				log.Debugf("Connection local address: %s, remote address: %s", conn.LocalAddr(), conn.RemoteAddr())
			}
			return conn, err
		}
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to MongoDB: %v", err)
	}

	//TODO: Look into modes and set here

	return session.DB(cfg.Name), nil
}

func (b *Backends) IsConnected() bool {
	return b.connected
}
