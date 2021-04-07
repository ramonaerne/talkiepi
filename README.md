# talkiepi
![assembled1](doc/talkiepi_assembled_1.jpg "Assembled talkiepi 1")
![assembled2](doc/talkiepi_assembled_2.jpg "Assembled talkiepi 2")

talkiepi is a headless capable Mumble client written in Go, written for walkie talkie style use on the Pi using GPIO pins for push to talk and LED status.  It is a fork of [barnard](https://github.com/layeh/barnard), which was a great jump off point for me to learn golang and have something that worked relatively quickly.


## 3D printable enclosure

In the stl directory are the stl files for the enclosure I have designed specifically for the Raspberry Pi B+ board layout (I am using a Raspberry Pi 3 Model B) and the PCB and components from the [US Robotics USB Speakerphone](https://www.amazon.com/USRobotics-USB-Internet-Speakerphone-USR9610/dp/B000E6IL10/ref=sr_1_1?ie=UTF8&qid=1472691020&sr=8-1&keywords=us+robotics+speakerphone).
I will be posting a blog post shortly with a full component list and build guide.  For more information regarding building a full talkiepi device, go check out my blog at [projectable.me](http://projectable.me).


## Installing talkiepi
Run the following command as the root user
```
su -i
./update-talkiepi.sh
```
which does the same as the guide below.

I have put together an install guide [here](doc/README.md).


## GPIO

```
J8:                        Connect to the phone circuit:
   3V3  (1) (2)  5V        
 GPIO2  (3) (4)  5V    
 GPIO3  (5) (6)  GND   
 GPIO4  (7) (8)  GPIO14
   GND  (9) (10) GPIO15
GPIO17 (11) (12) GPIO18    -> ringPwm 
GPIO27 (13) (14) GND       -> GND to ring circuit
GPIO22 (15) (16) GPIO23    -> ringEnable
   3V3 (17) (18) GPIO24
GPIO10 (19) (20) GND       -> GND to all button circuits below
 GPIO9 (21) (22) GPIO25    -> dialpin
GPIO11 (23) (24) GPIO8     -> pickup
   GND (25) (26) GPIO7     -> dialStartStop
 GPIO0 (27) (28) GPIO1 
 GPIO5 (29) (30) GND   
 GPIO6 (31) (32) GPIO12
GPIO13 (33) (34) GND   
GPIO19 (35) (36) GPIO16
GPIO26 (37) (38) GPIO20
   GND (39) (40) GPIO21
```

You can edit your pin assignments in `talkiepi.go`
```go
const (
	RingEnablePin   uint = 23
	RingPwmChannel  int  = 0 // pin 18
	...
)
```
and `gpio.go`
```go
var dialPin = pinDef{25, gpio.ActiveLow, gpio.EdgeRising, 0, []Event{EVENT_NOP, EVENT_DIAL_INC}}
var pickupPin = pinDef{8, gpio.ActiveLow, gpio.EdgeBoth, 0, []Event{EVENT_PICKUP_STOP, EVENT_PICKUP_START}}
var dialStartPin = pinDef{7, gpio.ActiveLow, gpio.EdgeBoth, 0, []Event{EVENT_DIAL_STOP, EVENT_DIAL_START}}

```

<!-- Here is a basic schematic of how I am currently controlling the LEDs and pushbutton:

![schematic](doc/gpio_diagram.png "GPIO Diagram") -->


## Pi Zero Fixes
I have compiled libopenal without ARM NEON support so that it works on the Pi Zero. The packages can be found in the [workarounds](/workarounds/). directory of this repo, install the libopenal1 package over your existing libopenal install.


## License

MPL 2.0

## Author

- talkiepi - [Daniel Chote](https://github.com/dchote)
- Barnard,Gumble Author - Tim Cooper (<tim.cooper@layeh.com>)

