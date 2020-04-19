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

type doorstate struct {
	ID        string    `json:"id"`
	State     string    `json:"state"`
	StateTime time.Time `json:"stateTime"`
}

// ServeHTTP implements the net/http Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusOK)
	accept := r.Header.Get("Accept")
	log.Printf("%#v\n", accept)
	switch r.Method {
	case "GET":
		state, err := s.ReadPin()
		if err != nil {
			state = fmt.Sprintf("%v", err)
		}
		if accept == "application/json" {
			//render json
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			ds := doorstate{ID: "Our door", State: state, StateTime: time.Now()}
			reply, err := json.MarshalIndent(&ds, "", "    ")
			if err != nil {
				reply = []byte(`{"error":"state marshal error"}`)
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
					"		<p>The door state is " + state + "</p></div>" +
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
