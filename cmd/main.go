package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/djerman3/gdo"
)

func main() {
	configfile := flag.String("config", "/etc/gdo/gdo.conf", "ets the config file path for secrets and ports and things.")
	testmode := flag.Bool("test", false, "Toggles true for test mode operation, where no PI controls actually work.")
	flag.Parse()

	cfg, err := gdo.NewConfig(*configfile)
	cfg.Testing = *testmode
	if err != nil {
		log.Println(err)
		return
	}
	s, err := gdo.NewWebserver(cfg)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%s listening on "+cfg.Server.Addr+":"+fmt.Sprintf("%d", cfg.Server.Port)+"\n", os.Args[0])
	log.Fatalln(s.ListenAndServe())
}
