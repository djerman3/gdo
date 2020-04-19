package gdapi

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/djerman3/gdo/piface"
	"github.com/luismesas/GoPi/spi"
)

// Server is the handler struct, carries the pi and mutex
type Server struct {
	rpi                  *piface.Digital
	piLock               sync.Mutex
	pinClosed            byte
	pinClosedAssertValue byte // add pinOpen assertions and "in motion" state if I ever get off my ass and install another switch
	relay                byte
}

// sensor pins assert one of these positions with one of their values
// when there's fewer pins than positions use induction, but assertions overrule

// Init : don't forget to init once
func (s *Server) Init() error {
	// set scan pins
	s.pinClosed = 5            //"this pin asserts door closed state"
	s.pinClosedAssertValue = 0 //"closed when zero"
	s.relay = 0
	// creates a new pifacedigital instance
	if s.rpi == nil {
		s.rpi = piface.NewDigital(spi.DEFAULT_HARDWARE_ADDR, spi.DEFAULT_BUS, spi.DEFAULT_CHIP)
		if s.rpi == nil {
			return fmt.Errorf("error on new rpi interface")
		}
	}
	return s.rpi.InitBoard()
}

//DoClick emulates a button click by cycling the relay 0.3 seconds
func (s *Server) DoClick() error {
	s.piLock.Lock()
	defer s.piLock.Unlock()
	s.rpi.Relays[s.relay].Toggle()
	time.Sleep(300 * time.Millisecond)
	s.rpi.Relays[s.relay].Toggle()
	return nil
}

//ReadPin emulates a button click by cycling the relay 0.3 seconds
func (s *Server) ReadPin() (string, error) {
	s.piLock.Lock()
	defer s.piLock.Unlock()
	reply := "Open"

	if s.rpi.InputPins[s.pinClosed].Value() == s.pinClosedAssertValue {
		reply = "Closed"
	}

	return reply, nil
}

//HTTP stuff

// ServeHTTP implements the net/http Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusOK)
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		state, err := s.ReadPin()
		if err != nil {
			state = fmt.Sprintf("%v", err)
		}
		w.Write([]byte(
			"<!doctype html> " +
				"<meta http-equiv=\"refresh\" content=\"10\">" +
				"		<title>Garage Door</title>" +
				"		<body>" +
				"		<h1>Garage Door</h1><p/>" +
				"		<h2>The door state is " + state + "</h2><p/>" +
				"       <form action=\"/\" method=\"POST\">" +
				"			<button type=\"submit\" formaction=\"/\" autofocus=\"autofocus\">CLICK</button>" +
				"			<button type=\"reset\" formaction=\"/\" >RELOAD</button>" +
				"		</form>" +
				"	</body>"))
	case "POST":
		w.Header().Set("Content-Type", "text/html")
		err := s.DoClick()
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "http://"+r.Host+"/", 301)

	}

}
