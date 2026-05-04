package aryan

import "encoding/json"

type ParamEntry struct {
	Name       string `json:"Name"`
	Value      any    `json:"Value,omitempty"`
	ArrayValue []any  `json:"Array_Value,omitempty"`
}

func (p ParamEntry) MarshalJSON() ([]byte, error) {
	if p.ArrayValue != nil {
		return json.Marshal(struct {
			Name       string `json:"Name"`
			ArrayValue []any  `json:"Array_Value"`
		}{Name: p.Name, ArrayValue: p.ArrayValue})
	}
	return json.Marshal(struct {
		Name  string `json:"Name"`
		Value any    `json:"Value"`
	}{Name: p.Name, Value: p.Value})
}

type ParamsPayload struct {
	ID     string       `json:"id"`
	Params []ParamEntry `json:"Params"`
}
