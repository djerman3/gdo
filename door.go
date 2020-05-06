// Package gdo provides the garage door webserver
// this file organizes the door part
package gdo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/djerman3/pimonitor"
)

//DoorState is the exposed door data, serializable
type DoorState struct {
	ID        string    `json:"id"`
	State     string    `json:"state"`
	StateTime time.Time `json:"stateTime"`
	Asserted  bool      `json:"asserted,omitempty"`
	Error     string    `json:"error,omitempty"`
}
type gdoDoor struct {
	Addr                 string
	PinClosed            int
	PinClosedAssertValue bool
	Relay                int
	Cmd                  map[string]string
	Client               *http.Client
	Testing              bool
}

func (d *gdoDoor) getJSON(url string, target interface{}) error {
	if d.Testing {
		// disable call effects for test mode
		return nil
	}
	r, err := d.Client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// Init : don't forget to init once
func (d *gdoDoor) Init(cfg *Config) error {
	// set scan pins
	d.PinClosed = cfg.Door.ClosedPin              // closedPin //"this pin asserts door closed state"
	d.PinClosedAssertValue = cfg.Door.ClosedValue //"closed when zero"
	d.Relay = cfg.Door.ClickRelay
	// creates a new pifacedigital instance
	d.Addr = cfg.Door.PiMonitorAddress
	d.Cmd = map[string]string{
		"click": fmt.Sprintf("http://%s/relay/%d/click", d.Addr, d.Relay),
		"read":  fmt.Sprintf("http://%s/input/%d/value", d.Addr, d.PinClosed),
	}
	d.Client = &http.Client{Timeout: 8 * time.Second}
	return nil
}

//DoClick emulates a button click by cycling the Relay 0.3 seconds
func (d *gdoDoor) DoClick() error {
	r := pimonitor.BoolValueResponse{}
	err := d.getJSON(d.Cmd["click"], &r)
	if err != nil {
		return err
	}
	if r.Error != "" {
		return fmt.Errorf("click Failed:%s", r.Error)
	}
	return nil
}

//ReadPin emulates a button click by cycling the Relay 0.3 seconds
func (d *gdoDoor) ReadPin() (string, bool, error) {
	r := pimonitor.BoolValueResponse{}
	err := d.getJSON(d.Cmd["read"], &r)
	if err != nil {
		return "", false, err
	}
	if r.Error != "" {
		return "", false, fmt.Errorf("click Failed:%s", r.Error)
	}
	// inferred value
	reply := "Open"
	// or do we assert "closed?"
	a := (r.Value == d.PinClosedAssertValue)
	if a {
		reply = "Closed" //asserted
	}

	return reply, a, nil
}
