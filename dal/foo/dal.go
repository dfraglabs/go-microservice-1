package foo

import (
	"context"
	"errors"
	"time"
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"

	"github.com/dfraglabs/go-microservice-1/dal/dalutil"
	"github.com/dfraglabs/go-microservice-1/dal/foo/client"
	"github.com/dfraglabs/go-microservice-1/dal/foo/types"
	"github.com/dfraglabs/go-microservice-1/deps/backends"
)

const (
	FOO_COLLECTION_NAME = "foo"
)

var log = logrus.WithField("pkg", "fdal")

//go:generate counterfeiter -o ../../fakes/foodal/dal.go . IDAL

type IDAL interface {
	GetBar(ctx context.Context, id int) (*types.Bar, error)
}

type DAL struct {
	indexes        []*mgo.Index
	fooClient    client.IClient

	*dalutil.SmartCollection
}

func NewFooDAL(be *backends.Backends, expiresAfterSec int) (*DAL, error) {
	fd := &DAL{
		indexes: []*mgo.Index{
			{
				Name:   "unique-field",
				Key:    []string{"foo-field"},
				Unique: true,
			},
			{
				Name:        "expiring-field",
				Key:         []string{"expires-after"},
				ExpireAfter: time.Second * time.Duration(expiresAfterSec),
			},
		},
	}

	if !be.IsConnected() {
		return nil, errors.New("DAL is not connected. Connect the parent DAL first")
	}

	fd.SmartCollection = dalutil.NewSmartCollection(be.MongoDB.C(FOO_COLLECTION_NAME), time.Minute)

	if err := fd.EnsureIndexes(fd.indexes); err != nil {
		return nil, err
	}

	fd.fooClient = be.FooClient

	return fd, nil
}

func (f *DAL) GetBar(ctx context.Context, id int) (*types.Bar, error) {
	// Fetch the data via a client
	bar, err := f.fooClient.GetBar(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get bar: %v", err)
	}

	// Do something with bar
	bar.Value = bar.Value + 1

	return bar, nil
}

// Meets the go-health.ICheckable interface
func (f *DAL) Status() (interface{}, error) {
	stat := map[string]interface{}{}

	if err := f.checkHealth(); err != nil {
		stat["status"] = err.Error()
		return stat, err
	}

	stat["status"] = "ok"

	return stat, nil
}

// Attempt to fetch any, single element
func (f *DAL) checkHealth() error {
	var r interface{}
	err := f.Collection().Find(nil).One(r)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
	}

	return err
}
