package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/luismesas/goPi/spi"
	"github.com/djerman3/gdo/doorstate"
	"github.com/djerman3/gdo/doorstate/piface"
)

func main() {
	runtime.GOMAXPROCS(2)

	pfd := piface.NewDigital(spi.DEFAULT_HARDWARE_ADDR, spi.DEFAULT_BUS, spi.DEFAULT_CHIP)
	err := pfd.InitBoard()
	if err != nil {
		log.Fatalln(err)
	}
	//debug hang
	fmt.Printf("Value %v \n", pfd.InputPins[5].Value())
	rcfg := doorstate.StateConfig{
		ScanInterval: "1.5s",
	}
	cfg := &rcfg
	//cfg.Sensors = make([]doorstate.ScanSensor)
	//cfg.Sensors = make([]doorstate.ScanSensor, 1)
	cfg.Sensors = append(cfg.Sensors, doorstate.ScanSensor{
		ID:             0,
		Port:           5,
		Position:       "Closed",
		State:          0,
		AssertPosition: "Closed",
		AssertState:    0,
	})
	//err = json.Unmarshal([]byte("{     \"sensors\": [		{			 \"id\": 0,			 \"port\": 5,			 \"position\": \"Closed\",			 \"state\": 0,			 \"scanTime\": \"0001-01-01T00:00:00Z\",			 \"lastState\": 0,			 \"lastScanTime\": \"0001-01-01T00:00:00Z\",			 \"assertPosition\": \"Closed\",			 \"assertState\": 0		}   ],   \"scanInterval\": 500000000,   \"clickPort\": 0,   \"alertSubs\": null}"), &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancelMon := context.WithCancel(context.Background())
	// defer cancelMon() // commented out to call explicitly below
	go cfg.Monitor(ctx, pfd)
	for i := 0; i < 20; i++ {
		time.Sleep(2 * time.Second)
		js, err := json.MarshalIndent(cfg, "", "     ")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(js))
		}

	}
	cancelMon()
	return
}
