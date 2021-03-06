package talkiepi

import (
	"crypto/tls"

	"github.com/brian-armstrong/gpio"
	"github.com/dchote/gumble/gumble"
	"github.com/dchote/gumble/gumbleopenal"
	"gobot.io/x/gobot/sysfs"
)

// Raspberry Pi GPIO pin assignments (CPU pin definitions)
const (
	RingEnablePin   uint = 23
	RingPwmChannel  int  = 0 // pin 18
	RingMessageCode      = "ring-on-your-phone-already"
	RingFreqNs           = 20000000 // 50Hz
	RingDurationSec      = 5
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

type Talkiepi struct {
	Config *gumble.Config
	Client *gumble.Client

	Address   string
	TLSConfig tls.Config

	ConnectAttempts uint

	Stream *gumbleopenal.Stream

	ChannelName    string
	IsConnected    bool
	IsTransmitting bool

	GPIOEnabled bool
	Button      gpio.Pin
	ButtonState uint

	EventQueue   chan Event
	CurrentState State
	DialCounter  int

	RingEnable    gpio.Pin
	RingPwm       *sysfs.PWMPin
	IsRingingChan chan struct{}

	AssignedNumber int
	NotReally      bool
}
