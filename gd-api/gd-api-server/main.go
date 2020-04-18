package main

import "log"

func main() {
	s := gdApi.Server{}
	err := s.Init()
	if err != nil {
		log.Fatalln(err)
	}

}
