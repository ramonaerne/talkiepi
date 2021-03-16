package talkiepi

import (
	"fmt"
	"time"
	"strconv"
)

func (b* Talkiepi) Transition(e Event) State {
	switch e {
	case EVENT_PICKUP_START:
		return STATE_CALL
	case EVENT_PICKUP_STOP:
		return STATE_IDLE
	case EVENT_DIAL_START:
		if b.CurrentState == STATE_CALL {
			b.DialCounter = 0
			return STATE_DIAL
		}
	case EVENT_DIAL_INC:
		b.DialCounter++
	case EVENT_DIAL_STOP:
		if b.CurrentState == STATE_DIAL {
			return STATE_CALL
		}
	case EVENT_RING_RECEIVE:
		if b.CurrentState == STATE_IDLE {
			return STATE_RING
		} else {
			return STATE_CALL
		}

	}
	return b.CurrentState
}


// function allowed to block until state changes
// is called asynchronously from event queue handler
func (b* Talkiepi) HandleState() {
	switch b.CurrentState  {
	case STATE_CALL:
		// TODO switch to channel and send ring cmd if a call was dialed and not switched
		channel := strconv.Itoa(b.DialCounter)
		fmt.Println("TODO: start transmitting if not already")
		fmt.Println("entered channel ", channel)
		// start transmitting probably instead of at launch
		//b.Connect()
		b.ChangeChannel(channel)
		b.SendMessage("test")
		//b.TransmitStart()
	case STATE_IDLE:
		// st
		//b.CleanUp()
		//b.TransmitStop()
		//b.ChangeChannel(string(ASSIGNED_NUMBER))
		fmt.Println("TODO: stop transmitting ")
	case STATE_RING:
		fmt.Print("ring: ")
		for b.CurrentState  == STATE_RING {
			// TODO turn on bell
			time.Sleep(1 * time.Second)
			fmt.Print(".")
		}
		// TODO turn off bell
		fmt.Println("done")
	}
}
