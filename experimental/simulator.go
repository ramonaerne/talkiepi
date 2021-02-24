package experimental

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ASSIGNED_NUMBER = 7
)

var counter = 0
var channel = ASSIGNED_NUMBER

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

func keyReader(c chan Event) {
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < 10; i++ {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		switch text[0] {
		case 'q':
			c <- EVENT_DIAL_START
		case 'w':
			c <- EVENT_DIAL_STOP
		case 'e':
			c <- EVENT_PICKUP_START
		case 'r':
			c <- EVENT_PICKUP_STOP
		case 't':
			split := strings.Split(text, " ")
			channel, _ = strconv.Atoi(split[1])
			c <- EVENT_RING_RECEIVE
		}
	}
	close(c)
}

type State int
const (
	STATE_IDLE = iota
	STATE_CALL
	STATE_DIAL
	STATE_RING
)

func (s State) String() string {
	return [...]string{"STATE_IDLE", "STATE_CALL", "STATE_DIAL", "STATE_RING"}[s]
}

func (s State) Transition(e Event) State {
	switch e {
	case EVENT_PICKUP_START:
		return STATE_CALL
	case EVENT_PICKUP_STOP:
		return STATE_IDLE
	case EVENT_DIAL_START:
		if s == STATE_CALL {
			return STATE_DIAL
		}
	case EVENT_DIAL_STOP:
		if s == STATE_DIAL {
			channel = counter
			return STATE_CALL
		}
	case EVENT_RING_RECEIVE:
		if s == STATE_IDLE {
			return STATE_RING
		} else {
			return STATE_CALL
		}

	}
	return s
}


// function allowed to block until state changes
// is called asynchronously from event queue handler
func (s* State) Handle() {
	switch *s {
	case STATE_DIAL:
		for *s == STATE_DIAL {
			// TODO turn on dialing interrupt logic here instead of using event queue
			time.Sleep(1 * time.Second)
			counter++
		}
		counter = 0
	case STATE_CALL:
		// TODO switch to channel and send ring cmd if a call was dialed and not switched
		fmt.Println("entered channel ", channel)
	case STATE_IDLE:
		// nop
	case STATE_RING:
		fmt.Print("ring: ")
		for *s == STATE_RING {
			// TODO turn on bell
			time.Sleep(1 * time.Second)
			fmt.Print(".")
		}
		// TODO turn off bell
		fmt.Println("done")
	}
}

func main() {
	fmt.Println("test")
	event_queue := make(chan Event)
	go keyReader(event_queue)

	s := State(STATE_IDLE)
	for i := range event_queue {
		fmt.Println(i)
		s = s.Transition(i)
		fmt.Println("transition to state:", s)
		go s.Handle()
	}
}