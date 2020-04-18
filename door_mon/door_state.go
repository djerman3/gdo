// Package doorstate implements the GDO door state-machine
// The door state machine tracks the status of the sensor(s) and the garage door and reports last
// measured (or inferred) state, and sends alerts to subscribed endpoint handlers when something changes.
// emits alerts per configdb based on reading the door sensor
// featurs: piface sensor leads for settings
//          logic for one-sensor or two-sensor ops
//          rest-like http(s) controls
//
//          path - description
//          /door/{n}/sensor/{n}
//              GET - read state of sensor
//              POST  .../sensor/{n}?port=#number&{off|on}={open|closed}
//                  - set sensor n to read input #number with 0 (off)
//					or 1 (on) stting a positive assertion of state and
//					the other state implying not-the-positive assertion
//					(but not absolutely the negative)
//				DELETE - remove sensor {n}
//          /door/{n}/sensors
//              GET - read state of all sensors
//              POST  .../sensor/{n}?port=#number&{off|on}={open|closed}
//                  - set sensor definitions with json
//				DELETE - drop a sensor
//			/door/{n}/
//          /door/{n}/interval
//              GET - read scan interval
//              POST - set scan interval
//          /door/{n}
//				GET - returns state and last state
//				POST - creates door entity
//				DELETE - drops door from config
//				PUT - modifies door config attribute(s) per payload (missing elements unchanged except where obsolescend)
//
//
package doorstate

import (
	"context"
	"fmt"
	"sync"
	"time"

	"jerman.info/gdo/doorstate/piface"
)

// ScanSensor stores the sensor definition, current and previous scan states (if present)
//  and configures the positively-asserted state for an indicated toggle reading (if present)
//  * presence of an asserted state and asserting position implies other states are inferred
//  * and therefore uncertain.  The positively-asserted position strongly asserts the asserted state.
//  * no asserted position/state indicates we're just guessing the state based on prior performance.
type ScanSensor struct {
	ID             int32     `json:"id"`
	Port           byte      `json:"port,omitEmpty"`
	Position       string    `json:"position,omitEmpty"`
	State          byte      `json:"state,omitEmpty"`
	ScanTime       time.Time `json:"scanTime,omitEmpty"`
	LastState      byte      `json:"lastState,omitEmpty"`
	LastScanTime   time.Time `json:"lastScanTime,omitEmpty"`
	AssertPosition string    `json:"assertPosition,omitEmpty"`
	AssertState    byte      `json:"assertState,omitEmpty"`
}

// AlertSubscription organizes the posting info for alerts including fail counter and last-good connect time for ageing and pruning
type AlertSubscription struct {
	ID       int32  `json:"id"`
	Name     string `json:"name,omitEmpty"`
	URLType  int32  `json:"urlType,omitEmpty"`
	URLText  string `json:"urlText,omitEmpty"`
	URLAuth  string `json:"urlAuth,omitEmpty"`
	URLKey64 string `json:"urlKey64,omitEmpty"`
}

// StateConfig organizes current configuration and readings for the state-machine
type StateConfig struct {
	Sensors      []ScanSensor        `json:"sensors"`
	ScanInterval string              `json:"scanInterval"`
	ClickPort    int8                `json:"clickPort"`
	AlertSubs    []AlertSubscription `json:"alertSubs"`
	Padlock      sync.Mutex          `json:"-"` //do not serialize  the mutex
}

// Poll reads the sensor and moves-back the previous reading
func (ss *ScanSensor) Poll(ctx context.Context, pfd *piface.Digital) error {
	if ss == nil {
		return fmt.Errorf("nil scansensor pointer")
	}
	// debug
	fmt.Println("Scanninig")
	ss.LastScanTime = ss.ScanTime
	ss.LastState = ss.State
	//read switch
	ss.ScanTime = time.Now()
	ss.State = pfd.InputPins[ss.Port].Value()
	//indicator := piface.ReadInput(ss.Port)
	if len(ss.AssertPosition) <= 0 {
		ss.AssertPosition = "Closed"
	}
	//Interpret state
	if ss.State == ss.AssertState {
		ss.Position = ss.AssertPosition
	} else {
		if ss.AssertPosition == "Opem" {
			ss.Position = "Closed"
		} else {
			ss.Position = "Open"
		}
	}

	return nil
}

// Poll reads the sensors in the collection to update their states
// This is where the mutex is used - outside the loop as there will be just a few sensors
func (s *StateConfig) Poll(ctx context.Context, pfd *piface.Digital) error {
	// scan over sensors
	fmt.Println("Polling")
	s.Padlock.Lock()
	defer s.Padlock.Unlock()
	for _, sswitch := range s.Sensors {
		err := sswitch.Poll(ctx, pfd)
		if err != nil {
			return err
		}
	}

	return nil
}

// Monitor runs until cancelled (or otherwise forcibly interrupted)
// It safely updates the stateconfig sensors and signals the alerts when a change in state is detected.
func (s *StateConfig) Monitor(ctx context.Context, pfd *piface.Digital) error {
	if s == nil {
		return fmt.Errorf("nil stateconfig pointer")
	}
	if len(s.ScanInterval) == 0 {
		s.ScanInterval = "5s"
	}
	interval, err := time.ParseDuration(s.ScanInterval)

	// coerce on err
	if err != nil {
		s.ScanInterval = "5s"
		interval, err = time.ParseDuration(s.ScanInterval)
		if err != nil {
			return err
		}
	}

	for {

		// sleep
		time.Sleep(interval)
		// are we there yet?
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// scan over sensors
			err := s.Poll(ctx, pfd)
			if err != nil {
				return err
			}
		}
	}
}
