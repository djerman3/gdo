package gdapi

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/djerman3/gdo/piface"
	"github.com/luismesas/GoPi/spi"
)

// Server is the handler struct, carries the pi and mutex
type Server struct {
	rpi    *piface.Digital
	piLock sync.Mutex
	pins   []byte
}

// Init : don't forget to init once
func (s *Server) Init() error {
	// set scan pins
	s.pins = []byte{0, 1, 2, 3, 4, 5, 6, 7}
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
	s.rpi.Relays[0].Toggle()
	time.Sleep(300 * time.Millisecond)
	s.rpi.Relays[0].Toggle()
	return nil
}

//ReadPin emulates a button click by cycling the relay 0.3 seconds
func (s *Server) ReadPin() (string, error) {
	s.piLock.Lock()
	defer s.piLock.Unlock()
	reply := "["
	for _, v := range s.pins {
		if len(reply) > 2 {
			reply += ", "
		}
		value := s.rpi.Relays[v].Value()
		reply += fmt.Sprintf("%d", int(value))
	}
	reply += "]"
	return reply, nil
}

//HTTP stuff

// ServeHTTP implements the net/http Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		state, err := s.ReadPin()
		if err != nil {
			state = fmt.Sprintf("%v", err)
		}
		w.Write([]byte(`{"message": "get called", "state":` + state + `}`))
	case "POST":
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "post called"}`))
	}

}
