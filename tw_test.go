package tw

import (
	"math/rand"
	"testing"
	"time"
)

func set(tw *TimerWheel, secoundCount int, counter *int) (err error) {
	setTime := time.Now().Add(time.Second * time.Duration(rand.Intn(secoundCount-1)+1))
	setValue := "test"
	err = tw.Set(setTime, setValue)
	*counter++
	return
}

func create(secoundCount int, counter *int) TimerWheel {
	return TimerWheel{
		Duration:     time.Duration(secoundCount) * time.Second,
		MaxSlotEvent: maxSlotEvent,
		CallBack: func(value string) error {
			*counter--
			return nil
		},
	}
}

var maxSlotEvent = 10

func TestTotal_Success(t *testing.T) {
	counter := 0
	secoundCount := 10
	twHeel := create(secoundCount, &counter)

	err := twHeel.Init()
	if err != nil {
		t.Error("error on init", err)
	}
	for i := 0; i < 10; i++ {
		set(&twHeel, secoundCount, &counter)
	}

	time.Sleep(time.Duration(secoundCount) * time.Second)
	if counter != 0 {
		t.Error("GetConfig error on reading data")
	}
}

func TestSet_Limit(t *testing.T) {
	counter := 0
	secoundCount := 5
	twHeel := create(secoundCount, &counter)

	err := twHeel.Init()
	if err != nil {
		t.Error("error on init", err)
	}
	for i := 0; i < maxSlotEvent*secoundCount+1; i++ {
		err := set(&twHeel, secoundCount, &counter)
		if err != nil {
			return
		}
	}
	t.Error("Don't say about error 'All slot is full'")
}
