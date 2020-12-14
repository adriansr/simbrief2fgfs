package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	sb2fgfs "github.com/adriansr/simbrief2fgfs/go-sb2fgfs"
)

const (
	MaxInputSize = 256 * 1024
)

func writeError(w http.ResponseWriter, err error, keysAndValues... string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	jsErr := map[string]interface{} {
		"message": err.Error(),
	}
	jsBody := map[string]interface{} {
		"error": jsErr,
	}
	for i := 0; i+1 < len(keysAndValues); i+=2 {
		jsErr[keysAndValues[i]] = keysAndValues[i+1]
	}
	body, err := json.Marshal(jsBody)
	if err != nil {
		body = []byte(`"error":{"message":"json marshal failed on error"}`)
	}
	w.Write(body)
}

func ConvertFP(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == -1 {
		w.WriteHeader(http.StatusLengthRequired)
		return
	}
	if r.ContentLength > MaxInputSize {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	source, err := ioutil.ReadAll(io.LimitReader(r.Body, r.ContentLength))
	r.Body.Close()
	if err != nil && err != io.EOF {
		writeError(w, err, "at", "input")
		return
	}
	converted, err := sb2fgfs.ConvertFromBytes(source)
	if err != nil {
		writeError(w, err, "at", "convert")
	}
	w.Header().Add("Content-Type", "application/xml")
	fmt.Fprint(w, string(converted))
}