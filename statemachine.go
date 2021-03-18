package talkiepi

import (
	"fmt"
	"time"
	"strconv"
)

func dialToChannel(counter int) string {
	// TODO find proper conversion if needed
	return strconv.Itoa(counter)
}

func (b* Talkiepi) Transition(e Event) State {
	switch e {
	case EVENT_PICKUP_START:
		if b.CurrentState == STATE_RING {
			close(b.IsRingingChan)
		}
		b.DialCounter = ASSIGNED_NUMBER
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
			if b.NotReally {
				fmt.Println("dialing: ", dialToChannel(b.DialCounter))
			}
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
		channel := dialToChannel(b.DialCounter)
		fmt.Println("TODO: start transmitting if not already")
		fmt.Println("entered channel ", channel)
		// start transmitting probably instead of at launch
		//b.Connect()
		b.ChangeChannel(channel)
		time.Sleep(time.Millisecond * 500)
		b.SendMessage(RING_MESSAGE_CODE)
		//b.TransmitStart()
	case STATE_IDLE:
		// st
		//b.CleanUp()
		b.TransmitStop()
		b.ChangeChannel("ring")
		//b.ChangeChannel(string(ASSIGNED_NUMBER))
		fmt.Println("TODO: stop transmitting ")
	case STATE_RING:
		b.IsRingingChan = make(chan struct{})
		go func(w chan struct{}) {
			fmt.Print("ring: ")
			b.RingEnable.High()
			b.RingPwm.Enable(true)

			select {
			case <-w:
			case <-time.After(RING_DURATION_SEC * time.Second):
			}
			b.RingEnable.Low()
			b.RingPwm.Enable(false)
			fmt.Println("stopped ringing")
		}(b.IsRingingChan)
	}
}
