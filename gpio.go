package talkiepi

import (
	"fmt"
	"github.com/brian-armstrong/gpio"
	"github.com/stianeikeland/go-rpio"
	"gobot.io/x/gobot/sysfs"
)

type pinDef struct {
	pin        uint
	logicLevel gpio.LogicLevel
	eventEdge  gpio.Edge
	lastValue  uint
	event      []Event
}

// gnd, 25, 8, 7 is order on raspi gpio header
var dialPin = pinDef{25, gpio.ActiveLow, gpio.EdgeRising, 0, []Event{EVENT_NOP, EVENT_DIAL_INC}}
var pickupPin = pinDef{8, gpio.ActiveLow, gpio.EdgeBoth, 0, []Event{EVENT_PICKUP_STOP, EVENT_PICKUP_START}}
var dialStartPin = pinDef{7, gpio.ActiveLow, gpio.EdgeBoth, 0, []Event{EVENT_DIAL_STOP, EVENT_DIAL_START}}

func (b *Talkiepi) initGPIO() {
	// we need to pull in rpio to pullup our button pin
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		b.GPIOEnabled = false
		return
	} else {
		b.GPIOEnabled = true
	}

	// setup pins
	pinCollection := []pinDef{dialPin, pickupPin, dialStartPin}
	m := make(map[uint]pinDef)

	// set pullup via rpio package
	// TODO: this fails once for every pin after boot
	for _, p := range pinCollection {
		m[p.pin] = p
		rpio.PullMode(rpio.Pin(p.pin), rpio.PullUp)
	}
	rpio.Close()

	//      var wg sync.WaitGroup
	//      defer wg.Done()

	watcher := gpio.NewWatcher()
	for _, p := range pinCollection {
		watcher.AddPinWithEdgeAndLogic(p.pin, gpio.EdgeBoth, p.logicLevel)
	}

	b.EventQueue = make(chan Event)
	go listenToInput(b.EventQueue, watcher, m)

	// then we can do our gpio stuff
	b.RingEnable = gpio.NewOutput(RingEnablePin, false)
	// setup pwm pin
	b.RingPwm = sysfs.NewPWMPin(RingPwmChannel)
	err := b.RingPwm.Export()
	if err != nil {
		fmt.Println("pwm pin export failed")
		b.GPIOEnabled = false
		return
	}
	b.RingPwm.SetPeriod(RingFreqNs) // [ns]
	if err != nil {
		fmt.Println("pwm SetPeriod failed")
		b.GPIOEnabled = false
		return
	}
	b.RingPwm.SetDutyCycle(RingFreqNs / 2) // 50% [ns]
	if err != nil {
		fmt.Println("pwm SetDutyCycle failed")
		b.GPIOEnabled = false
		return
	}
}

func (b *Talkiepi) LEDOn(LED gpio.Pin) {
	if b.GPIOEnabled == false {
		return
	}

	LED.High()
}

func (b *Talkiepi) LEDOff(LED gpio.Pin) {
	if b.GPIOEnabled == false {
		return
	}

	LED.Low()
}

func (b *Talkiepi) LEDOffAll() {
	if b.GPIOEnabled == false {
		return
	}

	b.LEDOff(b.OnlineLED)
	b.LEDOff(b.ParticipantsLED)
	b.LEDOff(b.TransmitLED)
}

func listenToInput(eventQueue chan Event, watcher *gpio.Watcher, m map[uint]pinDef) {
	for {
		pinNum, value := watcher.Watch()
		pin, present := m[pinNum]

		// skip unwatched events, this shouldn't happen
		if !present {
			continue
		}

		// debounce logic, trigger watcher on both edges
		// but only trigger events when value toggled (and event edge is desired)
		if pin.lastValue == value {
			continue
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

		//fmt.Printf("read %d from gpio %d\n", value, pin)
		//fmt.Println("event: ", pin.event[value])
		eventQueue <- pin.event[value]
	}
}
