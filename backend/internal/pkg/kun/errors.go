package kun

import (
	"errors"
	"fmt"
)

type KUNError struct {
	Code    string
	Message string
	Errors  []string
}

func (e *KUNError) Error() string {
	return fmt.Sprintf("KUN API error: code=%s, message=%s", e.Code, e.Message)
}

func IsKUNError(err error) (*KUNError, bool) {
	var kunErr *KUNError
	if errors.As(err, &kunErr) {
		return kunErr, true
	}
	return nil, false
}
