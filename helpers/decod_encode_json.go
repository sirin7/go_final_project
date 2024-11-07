package helpers

import (
	"encoding/json"
	"io"
)

// Десериализация и сериализация JSON
func DecodeJSON(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	decoder := json.NewDecoder(body)
	return decoder.Decode(v)
}

func EncodeJSON(w io.Writer, v interface{}) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(v)

}
