package main

import (
	"github.com/adriansr/simbrief2fgfs/go-sb2fgfs/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/convertfp", api.ConvertFP)
	log.Fatal(http.ListenAndServe(":8383", nil))
}
