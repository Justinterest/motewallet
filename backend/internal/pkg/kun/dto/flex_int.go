package dto

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexInt unmarshals JSON numbers or numeric strings from KUN APIs.
type FlexInt int

func (f *FlexInt) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*f = 0
		return nil
	}

	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexInt(n)
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if s == "" {
			*f = 0
			return nil
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("flex int: parse %q: %w", s, err)
		}
		*f = FlexInt(v)
		return nil
	}

	return fmt.Errorf("flex int: unsupported json value %s", string(data))
}

func (f FlexInt) Int() int {
	return int(f)
}
