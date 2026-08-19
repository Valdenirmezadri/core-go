package helperecho

import "encoding/json"

type Result struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (Result) Unmarshall(b []byte) (v Result, err error) {
	return _unmarshall[Result](b)
}

type Err struct {
	Errors []string `json:"errors"`
}

func (Err) Unmarshall(b []byte) (v Err, err error) {
	return _unmarshall[Err](b)
}

func _unmarshall[T Result | Err](b []byte) (v T, err error) {
	return v, json.Unmarshal(b, &v)
}
