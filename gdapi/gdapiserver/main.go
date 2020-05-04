package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/djerman3/gdo/gdapi"
)

func main() {
	port := flag.Int("port", 5000, "the Listener port for the service.")
	host := flag.String("host", "clubhouse.jerman.info", "The DNS server name for the service host.")
	addr := flag.String("addr", "0.0.0.0", "The IP listener address for the service.")
	testmode := flag.Bool("test", false, "Toggles true for test mode operation, where no PI controls actually work.")
	flag.Parse()

	s := &gdapi.Server{}
	err := s.Init(host, addr, port, testmode)
	if err != nil {
		log.Println(err)
		if err, ok := err.(*gdapi.TestModeError); ok {
			log.Fatalf("test mode unexpected:%v\n", err)
		}
	}
	log.Printf("%s listening on "+s.Address+":"+fmt.Sprintf("%d", s.Port)+"\n", os.Args[0])
	log.Fatalln(http.ListenAndServe(s.Address+":"+fmt.Sprintf("%d", s.Port), s))
}
