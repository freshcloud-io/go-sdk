package fresh

import (
	"fmt"
	"net/rpc"

	"github.com/freshcloud-io/go-sdk/utils"
	pb "github.com/freshcloud-io/protos/go/freshcloud"
	"github.com/hashicorp/go-plugin"
)

const (
	PluginNameCron       string = "fresh-cron"
	MagicCookieKeyCron   string = "FRESH_CRON_PLUGIN"
	MagicCookieValueCron string = ""
)

func NewFresherCron(f FreshCron) error {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   MagicCookieKeyCron,
			MagicCookieValue: MagicCookieValueCron,
		},
		Plugins: map[string]plugin.Plugin{
			PluginNameCron: &FreshCronPlugin{Impl: f},
		},
	})
	return nil
}

type FreshCron interface {
	Execute(args *pb.CronExecuteRequest) *pb.CronExecuteResponse
	GetCron() *Cron
}

type FreshCronRPC struct {
	client *rpc.Client
}

func (f *FreshCronRPC) Execute(args *pb.CronExecuteRequest) *pb.CronExecuteResponse {
	var resp *pb.CronExecuteResponse
	err := f.client.Call("Plugin.Execute", args, &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func (f *FreshCronRPC) GetCron() *Cron {
	var resp *Cron
	err := f.client.Call("Plugin.GetCron", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

type FreshCronRPCServer struct {
	Impl FreshCron
}

func (s *FreshCronRPCServer) Execute(args interface{}, resp *pb.CronExecuteResponse) error {
	a, ok := args.(*pb.CronExecuteRequest)
	if !ok {
		return fmt.Errorf("wrong cron arguments conversion")
	}

	*resp = *s.Impl.Execute(a)
	return nil
}

func (s *FreshCronRPCServer) GetCron(args interface{}, resp *Cron) error {
	*resp = *s.Impl.GetCron()
	return nil
}

type FreshCronPlugin struct {
	Impl FreshCron
}

func (p *FreshCronPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &FreshCronRPCServer{Impl: p.Impl}, nil
}

func (FreshCronPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &FreshCronRPC{client: c}, nil
}

type CronOption func(*Cron)

func WithCronName(name string) CronOption {
	return func(c *Cron) {
		c.Name = name
	}
}

func WithCronImage(image *Image) CronOption {
	return func(c *Cron) {
		c.Image = image
	}
}

func WithCronImageRegistry(registry *ImageRegistry) CronOption {
	return func(c *Cron) {
		c.ImageRegistry = registry
	}
}

func WithCronResource(res *Resources) CronOption {
	return func(c *Cron) {
		c.Resources = res
	}
}

type Cron struct {
	Name          string
	Image         *Image
	ImageRegistry *ImageRegistry
	Resources     *Resources
	SecretHandler *SecretHandler
	Dictionary    *Dictionary
	Queue         *Queue
	Volume        *Volume
}

func NewCron(opts ...CronOption) (*Cron, error) {
	var (
		resources = NewResource()
		name      = utils.RandomContainerName()
	)

	c := &Cron{
		Name:      name,
		Resources: resources,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *Cron) String() string {
	res := c.Resources.String()
	im := c.Image.String()
	imReg := c.ImageRegistry.String()

	return fmt.Sprintf("name: %s, resources: %s, image: %s, imageRegistry: %s", c.Name, res, im, imReg)
}
