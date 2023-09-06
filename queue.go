package fresh

import (
	"context"
	"fmt"
	"time"

	"github.com/freshcloud-io/go-sdk/grpc"
	pb "github.com/freshcloud-io/protos/go/freshcloud"
)

type Queue struct {
	client     pb.FreshcloudClient
	connection *grpc.Connection
}

func NewQueue() (*Queue, error) {
	// TODO: get these from somewhere
	conn, err := grpc.NewConnection("127.0.0.1:8401", 8400)
	if err != nil {
		return nil, err
	}
	return &Queue{
		client:     pb.NewFreshcloudClient(conn),
		connection: conn,
	}, nil
}

func (q *Queue) Produce(topic string, data []byte) (bool, error) {
	resp, err := q.client.ProduceValueQueue(context.Background(), &pb.ProduceValueQueueRequest{
		Topic:   topic,
		Message: data,
	})
	if err != nil {
		return false, fmt.Errorf("%s: %s", resp.Error, err.Error())
	}
	return resp.Ack, nil
}

func (q *Queue) Consume(topic string, messages chan<- []byte) error {
	req := &pb.ConsumeValuesQueueRequest{Topic: topic}
	stream, err := q.client.ConsumeValuesQueue(context.Background(), req)
	if err != nil {
		return err
	}
	stop := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-stop.C:
			// Tell the Server to close this Stream, used to clean up running on the server
			if err := stream.CloseSend(); err != nil {
				return err
			}
			close(messages)
			return nil
		default:
			// Receive on the stream
			res, err := stream.Recv()
			if err != nil {
				return err
			}
			messages <- res.Message
		}
	}
}
