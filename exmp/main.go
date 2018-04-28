package main

import (
	"fmt"
	"github.com/B2BFamily/tw"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	counter := 0
	secoundCount := 10
	twHeel := tw.TimerWheel{
		Duration:     time.Duration(secoundCount) * time.Second,
		MaxSlotEvent: 10,
		CallBack: func(value string) error {
			fmt.Printf("get %v value: %v\n", time.Now().Format("Mon Jan _2 15:04:05 2006"), value)
			counter--
			return nil
		},
	}

	go twHeel.Init()
	fmt.Println("tw init")
	for i := 0; i < 10; i++ {
		setTime := time.Now().Add(time.Second * time.Duration(rand.Intn(secoundCount)))
		setValue := strconv.Itoa(i)
		err := twHeel.Set(setTime, setValue)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("set %v in %v\n", setTime.Format("Mon Jan _2 15:04:05 2006"), setValue)
		counter++
	}

	time.Sleep(time.Duration(secoundCount) * time.Second)
	fmt.Printf("ready: %v\n", counter)
	var input string
	fmt.Scanln(&input)
}
