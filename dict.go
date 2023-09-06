package fresh

import (
	"context"
	"fmt"

	"github.com/freshcloud-io/go-sdk/grpc"
	pb "github.com/freshcloud-io/protos/go/freshcloud"
)

type Dictionary struct {
	client     pb.FreshcloudClient
	connection *grpc.Connection
}

func NewDictionary() (*Dictionary, error) {
	// TODO: get these from somewhere
	conn, err := grpc.NewConnection("127.0.0.1:8401", 8400)
	if err != nil {
		return nil, err
	}
	return &Dictionary{
		client:     pb.NewFreshcloudClient(conn),
		connection: conn,
	}, nil
}

func (d *Dictionary) Get(key string) ([]byte, error) {
	resp, err := d.client.GetValueDictionary(context.Background(), &pb.GetValueDictionaryRequest{
		Key: key,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %s", resp.Error, err.Error())
	}
	return resp.Value, nil
}

func (d *Dictionary) Exists(key string) (bool, error) {
	resp, err := d.client.ExistsValueDictionary(context.Background(), &pb.ExistsValueDictionaryRequest{
		Key: key,
	})
	if err != nil {
		return false, err
	}
	return resp.Exists, nil
}

// Put add the value to the list of other objects for the given key
// example: "key": ["val1","val2","val3"]
func (d *Dictionary) Put(key string, value []byte) error {
	resp, err := d.client.PutValueDictionary(context.Background(), &pb.PutValueDictionaryRequest{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("%s: %s", resp.Error, err.Error())
	}
	return nil
}

func (d *Dictionary) Delete(key string) error {
	resp, err := d.client.DeleteValueDictionary(context.Background(), &pb.DeleteValueDictionaryRequest{
		Key: key,
	})
	if err != nil {
		return fmt.Errorf("%s: %s", resp.Error, err.Error())
	}
	return nil
}

func (d *Dictionary) Pop(key string) ([]byte, error) {
	resp, err := d.client.PopValueDictionary(context.Background(), &pb.PopValueDictionaryRequest{
		Key: key,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %s", resp.Error, err.Error())
	}
	return resp.Value, nil
}

func (d *Dictionary) Length(key string) (int32, error) {
	resp, err := d.client.LengthDictionary(context.Background(), &pb.LengthDictionaryRequest{
		Key: key,
	})
	if err != nil {
		return -1, fmt.Errorf("%s: %s", resp.Error, err.Error())
	}
	return resp.Value, nil
}
