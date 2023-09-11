package fresh

import (
	"fmt"

	"github.com/freshcloud-io/go-sdk/utils"
)

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
