package helperecho

import (
	"encoding/json"
	"reflect"
)

type Result struct {
	Message string   `json:"message,omitempty"`
	Data    any      `json:"data,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

func (r Result) MarshalJSON() ([]byte, error) {
	data := r.Data
	if data != nil {
		v := reflect.ValueOf(data)
		if v.Kind() == reflect.Slice && v.IsNil() {
			data = reflect.MakeSlice(v.Type(), 0, 0).Interface()
		}
	}
	return json.Marshal(&struct {
		Message string   `json:"message,omitempty"`
		Data    any      `json:"data,omitempty"`
		Errors  []string `json:"errors,omitempty"`
	}{
		Message: r.Message,
		Data:    data,
		Errors:  r.Errors,
	})
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
