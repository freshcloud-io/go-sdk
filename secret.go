package fresh

import (
	"context"
	"os"

	"github.com/freshcloud-io/go-sdk/grpc"
	pb "github.com/freshcloud-io/protos/go/freshcloud"

	"github.com/joho/godotenv"
)

type SecretHandlerOption func(*SecretHandler)

func WithSecretHandlerEnvFile(filename string) SecretHandlerOption {
	return func(s *SecretHandler) {
		err := godotenv.Load(filename)
		if err != nil {
			panic(err)
		}
	}
}

func WithSecretHandlerEnvMap(envMap map[string]string) SecretHandlerOption {
	return func(s *SecretHandler) {
		for e, s := range envMap {
			os.Setenv(e, s)
		}
	}
}

type SecretHandler struct {
	ApplicationID string
	client        pb.FreshcloudClient
	connection    *grpc.Connection
}

func NewSecretHandler(applicationID string, opts ...SecretHandlerOption) (*SecretHandler, error) {
	// TODO: get these from somewhere
	conn, err := grpc.NewConnection("127.0.0.1:8401", 8400)
	if err != nil {
		return nil, err
	}

	sh := &SecretHandler{
		ApplicationID: applicationID,
		client:        pb.NewFreshcloudClient(conn),
		connection:    conn,
	}

	// load from our own APIs - users can set up envs directly in the UI or CLI
	resp, err := sh.client.LoadSecretsFromAPI(context.Background(), &pb.LoadSecretsFromAPIRequest{
		ApplicationId: applicationID,
	})
	if err != nil {
		return nil, err
	}
	for key, val := range resp.Values {
		os.Setenv(key, val)
	}

	// load from options
	// warn: this override any secret loaded from the backend
	for _, opt := range opts {
		opt(sh)
	}
	return sh, nil
}

func (s *SecretHandler) Add(key string, secret string) error {
	// TODO: should I save it on the backend as well? or backed into the image is sufficient?
	os.Setenv(key, secret)
	return nil
}
