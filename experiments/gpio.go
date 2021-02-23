package main

import (
	"fmt"
//x	"sync"
//	"time"
	"github.com/brian-armstrong/gpio"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	// we need to pull in rpio to pullup our button pin
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		return
	}

	ButtonPin := 25
	ButtonPinPullUp := rpio.Pin(ButtonPin)
	ButtonPinPullUp.PullUp()
//	var ButtonState uint = 1
	rpio.Close()

	// unfortunately the gpio watcher stuff doesnt work for me in this context, so we have to poll the button instead
//	Button := gpio.NewInput(25)
	//	dial := 0

//	var wg sync.WaitGroup
//	defer wg.Done()

	// watcher test
	watcher := gpio.NewWatcher()
	//watcher.AddPin(25)
	watcher.AddPinWithEdgeAndLogic(25, gpio.EdgeFalling, gpio.ActiveHigh)
	defer watcher.Close()


    		for {
        		pin, value := watcher.Watch()
        		fmt.Printf("read %d from gpio %d\n", value, pin)
    		}


	// then we can do our gpio stuff
	//b.OnlineLED = gpio.NewOutput(OnlineLEDPin, false)
	//b.ParticipantsLED = gpio.NewOutput(ParticipantsLEDPin, false)
	//b.TransmitLED = gpio.NewOutput(TransmitLEDPin, false)
}
