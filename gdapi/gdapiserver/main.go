package main

import (
	"log"
	"net/http"

	"github.com/djerman3/gdo/gdapi"
)

func main() {
	s := &gdapi.Server{}
	err := s.Init()
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(http.ListenAndServe("0.0.0.0:5000", s))
}
