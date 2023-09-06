package fresh

import (
	"fmt"
	"strings"
)

type Image struct {
	Commands []string
}

func NewImage() {

}

func (i *Image) CopyLocalDir() error { return nil }

func (i *Image) CopyLocalFile() error { return nil }

func (i *Image) CopyMount() error { return nil }

func (i *Image) DockerfileCommands() error { return nil }

func (i *Image) RunCommands() error { return nil }

func (i *Image) String() string {
	return fmt.Sprintf("commands: %s", strings.Join(i.Commands, ","))
}

type ImageRegistry struct {
	URL        string
	Repository string
}

func NewImageRegistry() {
}

func (ir *ImageRegistry) Resolve() error {
	return nil
}

func (ir *ImageRegistry) String() string {
	return fmt.Sprintf("url: %s\trepository: %s", ir.URL, ir.Repository)
}
