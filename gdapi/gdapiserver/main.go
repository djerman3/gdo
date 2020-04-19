package main

import (
	"log"
	"net/http"
	"os"

	"github.com/djerman3/gdo/gdapi"
)

func main() {
	s := &gdapi.Server{}
	err := s.Init()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.printf("%s listening on 0:5000\n", os.Args[0])
	log.Fatalln(http.ListenAndServe("0.0.0.0:5000", s))
}
