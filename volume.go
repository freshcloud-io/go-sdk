package fresh

import "fmt"

type VolumeOption func(*Volume)

func WithVolumeSize(size int) VolumeOption {
	return func(v *Volume) {
		v.Size = size
	}
}

type Volume struct {
	Name string
	Path string
	// Size is expressed in MB with a default of 150
	Size int
}

func NewVolume(name, path string, opts ...VolumeOption) (*Volume, error) {
	v := &Volume{
		Name: name,
		Path: path,
		Size: 150,
	}

	for _, opt := range opts {
		opt(v)
	}

	return v, nil
}

func (v *Volume) WalkDirectory() ([]string, error) {
	return nil, nil
}

func (v *Volume) String() string {
	return fmt.Sprintf("name: %s\tpath: %s\tsize: %d", v.Name, v.Path, v.Size)
}
