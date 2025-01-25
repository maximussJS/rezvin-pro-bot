package utils

import (
	"encoding/json"
	"net/http"
)

func MustWriteBytesToResponseWriter(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)

	PanicIfError(err)
}

func MustWriteJsonToResponseWriter(w http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(w).Encode(&data)

	PanicIfError(err)
}
