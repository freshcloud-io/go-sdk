package grpc

import (
	"time"

	"github.com/freshcloud-io/go-sdk/utils"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	retriableErrors = []codes.Code{codes.Unavailable, codes.DataLoss, codes.DeadlineExceeded}
	retryTimeout    = 50 * time.Millisecond
)

type Connection struct {
	*grpc.ClientConn
}

func NewConnection(bindAddr string, port int) (*Connection, error) {
	rpcAddr, err := utils.RPCAddr(bindAddr, port)
	if err != nil {
		return nil, err
	}

	unaryInterceptor := retry.UnaryClientInterceptor(
		retry.WithCodes(retriableErrors...),
		retry.WithMax(3),
		retry.WithBackoff(retry.BackoffLinear(retryTimeout)),
	)
	log.Debug().Msg("initialize gRPC client connection")
	conn, err := grpc.Dial(rpcAddr,
		grpc.WithBlock(), // this helps with having the grpc client waiting for the server to be avaialble and retry
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryInterceptor),
	)
	if err != nil {
		return nil, err
	}
	return &Connection{
		conn,
	}, nil
}

func (c *Connection) CloseConn() error {
	log.Debug().Msg("close gRCP client connection of Fresh Cloud")
	return c.Close()
}
