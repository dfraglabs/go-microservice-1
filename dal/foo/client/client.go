package client

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/dfraglabs/go-microservice-1/dal/foo/types"
)

//go:generate counterfeiter -o ../../../fakes/fooclient/client.go . IClient

var log = logrus.WithField("pkg", "fdal.client")

type IClient interface {
	GetBar(ctx context.Context, id int) (*types.Bar, error)
}

type Client struct {}

func NewFooClient(host, serviceName string) *Client {
	return &Client{}
}

func (t *Client) GetBar(ctx context.Context, id int) (*types.Bar, error) {
	return &types.Bar{}, nil
}

// Satisfy go-health.ICheckable interface
func (t *Client) Status() (interface{}, error) {
	return nil, nil
}