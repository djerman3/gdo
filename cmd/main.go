package main

import (
	"flag"
	"log"
	"os"

	"github.com/djerman3/gdo"
)

func main() {
	var configfile string
	var testmode bool
	flag.StringVar(&configfile, "config", "/etc/gdo/gdo.conf", "Sets the config file path for secrets and ports and things.")
	flag.BoolVar(&testmode, "test", false, "Toggles true for test mode operation, where no PI controls actually work.")
	flag.Parse()

	cfg, err := gdo.NewConfig(configfile)
	if err != nil {
		log.Println(err)
		return
	}
	cfg.Testing = testmode
	if err != nil {
		log.Println(err)
		return
	}
	s, err := gdo.NewWebserver(cfg)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%s starting\n", os.Args[0])
	log.Fatalln(s.ListenAndServe())
}
