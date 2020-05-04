package gdapi

import (
	"encoding/json"
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
	testmode             bool
	piLock               sync.Mutex
	PinClosed            byte   `json:"pinClosed"`
	PinClosedAssertValue byte   `json:"pinClosedAssertValue"` // add pinOpen assertions and "in motion" state if I ever get off my ass and install another switch
	Relay                byte   `json:"controlRelay"`
	Port                 int    `json:"port,omitempty"`
	Address              string `json:"address,omitempty"`
	Host                 string `json:"host,omitempty"`
}

// sensor pins assert one of these positions with one of their values
// when there's fewer pins than positions use induction, but assertions overrule

//TestModeError is an error to distinguish test mode, with a bool
// to indicate that test mode was requested (and the pi disabled on purpose)
// vs test mode fallback (Testmode = false) indicating the pi was not
// found/available
type TestModeError struct {
	err error
}

func (t *TestModeError) Error() string {
	return t.err.Error()
}

func (s *Server) initPi() error {

	if s.testmode {
		s.rpi = nil //force no pi
		return nil
	}
	if s.rpi == nil {
		s.rpi = piface.NewDigital(spi.DEFAULT_HARDWARE_ADDR, spi.DEFAULT_BUS, spi.DEFAULT_CHIP)
		if s.rpi == nil {
			err := fmt.Errorf("error on new rpi interface")
			t := TestModeError{
				err: err,
			}
			return &t
		}
	}
	err := s.rpi.InitBoard()
	if err != nil {
		t := TestModeError{
			err: err,
		}
		s.rpi = nil //dereference and carry on, server will be in test mode
		return &t
	}
	return nil
}

// Init : don't forget to init once
func (s *Server) Init(host *string, addr *string, port *int, testmode *bool) error {
	// set web stuff
	if port != nil {
		s.Port = *port
	}
	if host != nil {
		s.Host = *host
	}
	if addr != nil {
		s.Address = *addr
	}
	if testmode != nil {
		s.testmode = *testmode
	}
	// set scan pins
	s.PinClosed = 5            //"this pin asserts door closed state"
	s.PinClosedAssertValue = 0 //"closed when zero"
	s.Relay = 0
	// creates a new pifacedigital instance
	err := s.initPi()
	if err != nil {
		if err, ok := err.(*TestModeError); ok {
			s.testmode = true
			return err
		}
	}
	return err
}

//DoClick emulates a button click by cycling the Relay 0.3 seconds
func (s *Server) DoClick() error {
	s.piLock.Lock()
	defer s.piLock.Unlock()
	// guard
	if s.rpi != nil {
		s.rpi.Relays[s.Relay].Toggle()
		time.Sleep(300 * time.Millisecond)
		s.rpi.Relays[s.Relay].Toggle()
	} else {
		// do test mode
		log.Println("Test Mode: Click!")
	}
	return nil
}

//ReadPin emulates a button click by cycling the Relay 0.3 seconds
func (s *Server) readPin() (string, bool, error) {
	s.piLock.Lock()
	defer s.piLock.Unlock()
	// inferred value
	reply := "Open"
	a := false
	if s.rpi != nil {
		if s.rpi.InputPins[s.PinClosed].Value() == s.PinClosedAssertValue {
			// positive assertion
			reply = "Closed"
			a = true
		}
	} else {
		// do test mode
		log.Println("Test Mode: Read Open!")
	}

	return reply, a, nil
}

//HTTP stuff

//DoorState represents a snapshot of the garage door's status.  States are asserted or implied by sensors.
type DoorState struct {
	ID        string    `json:"id"`
	State     string    `json:"state"`
	StateTime time.Time `json:"stateTime"`
	Asserted  bool      `json:"asserted,omitempty"`
	Error     string    `json:"error,omitempty"`
}

//GetDoorState returns the state of the door
func (s *Server) GetDoorState() (DoorState, error) {
	ds := DoorState{ID: "Our door", StateTime: time.Now()}
	state, asserted, err := s.readPin()
	if err != nil {
		state = fmt.Sprintf("%v", err)
	}
	ds.State = state
	ds.Asserted = asserted
	return ds, err
}

// ServeHTTP implements the net/http Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusOK)
	accept := r.Header.Get("Accept")
	log.Printf("%#v\n", accept)
	// need this ds for all renderings
	ds, err := s.GetDoorState() // error included in struct
	if err != nil {
		ds.Error = err.Error()
	}
	switch r.Method {
	case "GET":
		if accept == "application/json" {
			//render json
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			reply, err := json.MarshalIndent(&ds, "", "    ")
			if err != nil {
				reply = []byte("{\"error\":" + err.Error() + "}")
			}
			w.Write(reply)

		} else {
			//render html
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(
				"<!doctype html> " + "<head>" +
					"<meta charset=\"utf-8\">" +
					"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">" +
					"<link rel=\"stylesheet\" href=\"https://maxcdn.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap.min.css\">" +
					"<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js\"></script>" +
					"<script src=\"https://maxcdn.bootstrapcdn.com/bootstrap/3.4.1/js/bootstrap.min.js\"></script>" +
					"<meta http-equiv=\"refresh\" content=\"6\">" +
					"		<title>Garage Door</title>" +
					"</head>" +
					"		<body>" +
					"		<div class=\"jumbotron text-center\"><h1>Garage Door</h1>" +
					"		<p>At " + ds.StateTime.Format("Jan 2 2006 15:04:05") + "</p>" +
					"       <p>The door state is " + ds.State + "</p></div>" +
					"       <form action=\"/\" method=\"POST\">" +
					"		 <div class=\"container\">" +
					"        <div class=\"row\"><div class=\"col-6\">" +
					"			<button class=\"btn btn-primary btn-lg btn-block\" type=\"submit\" formaction=\"/\" autofocus=\"autofocus\">CLICK</button>" +
					"        </div><div class=\"col-6\">" +
					"			<button class=\"btn btn-secondary btn-lg btn-block\" type=\"reset\" formaction=\"/\" >RELOAD</button>" +
					"		</div></div></div>" +
					"		</form>" +
					"	</body>"))
		}
	case "POST":
		w.Header().Set("Content-Type", "text/html")
		err := s.DoClick()
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "http://"+r.Host+"/", 301)

	}

}
