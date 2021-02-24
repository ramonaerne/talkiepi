package experimental

import (
	"fmt"
//x	"sync"
//	"time"
	"github.com/brian-armstrong/gpio"
	"github.com/stianeikeland/go-rpio"
)

type pinDef struct {
	pin uint
	logicLevel gpio.LogicLevel
	event []Event
}

// gnd, 25, 8, 7 is order on raspi gpio header
var dialPin      = pinDef{25,
	gpio.ActiveHigh,
	[]Event{EVENT_DIAL_INC, EVENT_NOP},
	}
var pickupPin    = pinDef{8, gpio.ActiveLow, []Event{EVENT_PICKUP_START, EVENT_PICKUP_STOP}}
var dialStartPin = pinDef{7, gpio.ActiveHigh, []Event{EVENT_DIAL_START, EVENT_DIAL_STOP}}
var pinCollection = []pinDef { dialPin, pickupPin, dialStartPin}

func watchPins(c chan Event) {
	// we need to pull in rpio to pullup our button pin
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		return
	}

	for _, p := range pinCollection {
		rpio.PullMode(rpio.Pin(p.pin), rpio.PullUp)
	}
	rpio.Close()

	//	var wg sync.WaitGroup
	//	defer wg.Done()

	watcher := gpio.NewWatcher()
	for _, p := range pinCollection {
		rpio.PullMode(rpio.Pin(p.pin), rpio.PullUp)
		watcher.AddPinWithEdgeAndLogic(p.pin, gpio.EdgeBoth, p.logicLevel)
	}
	defer watcher.Close()

	for {
		pin, value := watcher.Watch()
		for _, coll := range pinCollection {
			if coll.pin == pin {
				if coll.event[value] != EVENT_NOP {
					c <- coll.event[value]
				}
				break
			}
		}
		fmt.Printf("read %d from gpio %d\n", value, pin)
	}
}
