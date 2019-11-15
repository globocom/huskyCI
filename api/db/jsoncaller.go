package db

import (
	"encoding/json"
)

// Marshal will call Marshal function from json struct.
func (jCaller *JSONCaller) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal will call Unmarshal function from json struct.
func (jCaller *JSONCaller) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
