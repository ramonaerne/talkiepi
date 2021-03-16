package main 

import (
	"fmt"
//x	"sync"
//	"time"
	"github.com/brian-armstrong/gpio"
	"github.com/stianeikeland/go-rpio"
)

type Event int
const (
	EVENT_NOP = iota
	EVENT_PICKUP_START
	EVENT_PICKUP_STOP
	EVENT_DIAL_START
	EVENT_DIAL_STOP
	EVENT_RING_RECEIVE
	EVENT_DIAL_INC
)

func (e Event) String() string {
	return [...]string{"EVENT_NOP", "EVENT_PICKUP_START", "EVENT_PICKUP_STOP", "EVENT_DIAL_START", "EVENT_DIAL_STOP", "EVENT_RING_RECEIVE", "EVENT_DIAL_INC"}[e]
}

type pinDef struct {
	pin uint
	logicLevel gpio.LogicLevel
	eventEdge gpio.Edge
	lastValue uint
	event []Event
}

// gnd, 25, 8, 7 is order on raspi gpio header
var dialPin      = pinDef{25,gpio.ActiveLow, gpio.EdgeRising, 0, []Event{EVENT_NOP, EVENT_DIAL_INC}}
var pickupPin    = pinDef{8, gpio.ActiveHigh, gpio.EdgeBoth, 0, []Event{EVENT_PICKUP_STOP, EVENT_PICKUP_START}}
var dialStartPin = pinDef{7, gpio.ActiveLow, gpio.EdgeBoth, 0, []Event{EVENT_DIAL_STOP, EVENT_DIAL_START}}

func main() {
	// we need to pull in rpio to pullup our button pin
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		return
	}


	pinCollection := []pinDef { dialPin, pickupPin, dialStartPin}
	m := make(map[uint]pinDef)

	for _, p := range pinCollection {
		m[p.pin] = p
		rpio.PullMode(rpio.Pin(p.pin), rpio.PullUp)
	}
	fmt.Println("map: ", m)
	rpio.Close()

	//	var wg sync.WaitGroup
	//	defer wg.Done()

	watcher := gpio.NewWatcher()
	for _, p := range pinCollection {
		watcher.AddPinWithEdgeAndLogic(p.pin, gpio.EdgeBoth, p.logicLevel)
	}

	for {
		pinNum, value := watcher.Watch()
		pin, present := m[pinNum]

		// skip unwatched events, this shouldn't happen
		if !present {
			continue;
		}

		// debounce logic, trigger watcher on both edges
		// but only trigger events when value toggled (and event edge is desired)
		if pin.lastValue == value {
			 continue;
		}
		pin.lastValue = value
		m[pinNum] = pin

		// skip non-events
		switch {
		case pin.eventEdge == gpio.EdgeRising && value == 0:
			continue
		case pin.eventEdge == gpio.EdgeFalling && value == 1:
			continue
		}

		fmt.Printf("read %d from gpio %d\n", value, pin)
		fmt.Println("event: ", pin.event[value])
	}
}
