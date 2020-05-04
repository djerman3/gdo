// Package gdo provides the garage door webserver
// this file organizes the door part
package gdo

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/djerman3/gdo/piface"
	"github.com/luismesas/GoPi/spi"
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
	rpi                  *piface.Digital
	PinClosed            int
	PinClosedAssertValue int
	Relay                int
	piLock               sync.Mutex
}

func (d *gdoDoor) initPi() error {

	if d.rpi == nil {
		d.rpi = piface.NewDigital(spi.DEFAULT_HARDWARE_ADDR, spi.DEFAULT_BUS, spi.DEFAULT_CHIP)
		if d.rpi == nil {
			err := fmt.Errorf("error on new rpi interface")
			return err
		}
	}
	err := d.rpi.InitBoard()
	if err != nil {
		d.rpi = nil //dereference and carry on, server will be in test mode
		return err
	}
	return nil
}

// Init : don't forget to init once
func (d *gdoDoor) Init(cfg *Config) error {
	// set scan pins
	d.PinClosed = cfg.Door.ClosedPin              // closedPin //"this pin asserts door closed state"
	d.PinClosedAssertValue = cfg.Door.ClosedValue //"closed when zero"
	d.Relay = cfg.Door.ClickRelay
	// creates a new pifacedigital instance
	err := d.initPi()
	if err != nil {
		return err
	}
	return nil
}

//DoClick emulates a button click by cycling the Relay 0.3 seconds
func (d *gdoDoor) DoClick() error {
	d.piLock.Lock()
	defer d.piLock.Unlock()
	// guard
	if d.rpi != nil {
		d.rpi.Relays[d.Relay].Toggle()
		time.Sleep(300 * time.Millisecond)
		d.rpi.Relays[d.Relay].Toggle()
	} else {
		// do test mode
		log.Println("Test Mode: Click!")
	}
	return nil
}

//ReadPin emulates a button click by cycling the Relay 0.3 seconds
func (d *gdoDoor) readPin() (string, bool, error) {
	d.piLock.Lock()
	defer d.piLock.Unlock()
	// inferred value
	reply := "Open"
	a := false
	if d.rpi != nil {
		if d.rpi.InputPins[d.PinClosed].Value() == byte(d.PinClosedAssertValue) {
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
