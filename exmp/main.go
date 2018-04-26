package main

import (
	"fmt"
	"github.com/B2BFamily/tw"
	"time"
)

func main() {
	twHeel := tw.TimerWheel{
		SlotInterval: time.Duration(1) * time.Second,
		Duration:     time.Duration(5) * time.Second,
		MaxSlotEvent: 20,
	}

	go twHeel.Init()
	var input string
	fmt.Scanln(&input)
}
