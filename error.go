package fresh

import "encoding/gob"

type freshError struct {
	Message string `json:"message"`
}

func init() {
	gob.Register(&freshError{})
}

// New returns an error that formats as the given text.
func NewFreshError(text string) error {
	return &freshError{text}
}

func (e *freshError) Error() string {
	return e.Message
}

const version = 1

func (e *freshError) GobEncode() ([]byte, error) {
	r := make([]byte, 0, len(e.Message)+1)
	r = append(r, version)
	return append(r, e.Message...), nil
}

func (e *freshError) GobDecode(b []byte) error {
	if b[0] != version {
		return NewFreshError("gob decode of freshError failed: unsupported version")
	}
	e.Message = string(b[1:])
	return nil
}
