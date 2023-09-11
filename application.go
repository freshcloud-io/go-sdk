package fresh

import (
	"fmt"
	"net/rpc"

	"github.com/freshcloud-io/go-sdk/utils"
	"github.com/hashicorp/go-plugin"
)

const (
	PluginNameApplication       string = "fresh-application"
	MagicCookieKeyApplication   string = "FRESH_APPLICATION_PLUGIN"
	MagicCookieValueApplication string = "4ab23e44-f2c4-44c2-9b66-448d78d4713c"
)

func NewFresherApplication(f FreshApplication) error {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   MagicCookieKeyApplication,
			MagicCookieValue: MagicCookieValueApplication,
		},
		Plugins: map[string]plugin.Plugin{
			PluginNameApplication: &FreshApplicationPlugin{Impl: f},
		},
	})
	return nil
}

type FreshApplicationRPC struct {
	client *rpc.Client
}

func (f *FreshApplicationRPC) EntryPoint() error {
	var resp error
	err := f.client.Call("Plugin.EntryPoint", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func (f *FreshApplicationRPC) Shutdown() error {
	var resp error
	err := f.client.Call("Plugin.Shutdown", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func (f *FreshApplicationRPC) GetApplication() *Application {
	var resp *Application
	err := f.client.Call("Plugin.GetApplication", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

type FreshApplicationRPCServer struct {
	Impl FreshApplication
}

func (s *FreshApplicationRPCServer) EntryPoint(args interface{}, resp *error) error {
	*resp = s.Impl.EntryPoint()
	return nil
}

func (s *FreshApplicationRPCServer) Shutdown(args interface{}, resp *error) error {
	*resp = s.Impl.Shutdown()
	return nil
}

func (s *FreshApplicationRPCServer) GetApplication(args interface{}, resp *Application) error {
	*resp = *s.Impl.GetApplication()
	return nil
}

type FreshApplicationPlugin struct {
	Impl FreshApplication
}

func (p *FreshApplicationPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &FreshApplicationRPCServer{Impl: p.Impl}, nil
}

func (FreshApplicationPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &FreshApplicationRPC{client: c}, nil
}

type ApplicationOption func(*Application)

func WithApplicationName(name string) ApplicationOption {
	return func(a *Application) {
		a.Name = name
	}
}

func WithApplicationImage(image *Image) ApplicationOption {
	return func(a *Application) {
		a.Image = image
	}
}

func WithApplicationImageRegistry(registry *ImageRegistry) ApplicationOption {
	return func(a *Application) {
		a.ImageRegistry = registry
	}
}

func WithApplicationResource(res *Resources) ApplicationOption {
	return func(a *Application) {
		a.Resources = res
	}
}

type Application struct {
	Name          string
	Image         *Image
	ImageRegistry *ImageRegistry
	Resources     *Resources
}

func NewApplication(opts ...ApplicationOption) (*Application, error) {
	var (
		resources = NewResource()
		name      = utils.RandomContainerName()
	)

	a := &Application{
		Name:      name,
		Resources: resources,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a, nil
}

func (a *Application) String() string {
	res := a.Resources.String()
	im := a.Image.String()
	imReg := a.ImageRegistry.String()

	return fmt.Sprintf("name: %s, resources: %s, image: %s, imageRegistry: %s", a.Name, res, im, imReg)
}
