package fresh

import "fmt"

type FreshError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err *FreshError) Error() string {
	return fmt.Sprintf("[%s] %s", err.Code, err.Message)
}
