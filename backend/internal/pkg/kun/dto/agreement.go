package dto

import "encoding/json"

func (l *AgreementList) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*l = nil
		return nil
	}

	var items []Agreement
	if err := json.Unmarshal(data, &items); err == nil {
		*l = items
		return nil
	}

	var wrapped struct {
		List []Agreement `json:"list"`
	}
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}
	*l = wrapped.List
	return nil
}
