package fresh

import (
	"fmt"

	"github.com/freshcloud-io/go-sdk/utils"
)

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
