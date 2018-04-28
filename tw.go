package tw

import (
	"errors"
	"fmt"
	"time"
)

type TimerWheel struct {
	Duration     time.Duration      //Максимальное кол-во времени, которое будет обрабатываться
	MaxSlotEvent int                //максимальное кол-во событий на одну единицу времени
	CallBack     func(string) error //Обработчик событий
	slotInterval time.Duration      //Интервал времени обработки в секундах
	slots        [][]string         //Колесо, в котором будет все храниться
	ping         chan bool          //канал
	current      int                //текущий номер элемента
	slotCount    int                //количество слотов в колесе
}

func (base *TimerWheel) Init() error {
	if base.CallBack == nil {
		return errors.New("Callback is nil")
	}
	base.slotCount = int(base.Duration.Seconds())
	if base.slotCount < 1 {
		return errors.New("Duration less than a second")
	}
	base.current = 0
	base.slotInterval = time.Duration(1) * time.Second
	if base.MaxSlotEvent == 0 {
		base.MaxSlotEvent = 20
	}
	base.ping = make(chan bool)
	base.slots = make([][]string, base.slotCount)
	for i := 0; i < base.slotCount; i++ {
		base.slots[i] = make([]string, base.MaxSlotEvent)
	}
	go base.pinger()
	go base.handler()
	return nil
}

//Добавляем в TimerWheel событие на указанное время
func (base *TimerWheel) Set(timeEvent time.Time, valueEvent string) error {
	var duration = time.Duration(timeEvent.UTC().Unix()-time.Now().UTC().Unix()) * time.Second
	slotIndex := int(duration.Seconds() / base.slotInterval.Seconds())
	if slotIndex > base.slotCount {
		return errors.New(fmt.Sprintf("Time of event is too long, event time: %v, event duration: %v, max duration: %v", timeEvent, duration, base.Duration))
	}
	slotIndex = (slotIndex + base.current) % base.slotCount
	flag := false
	for i := slotIndex; i != slotIndex-1; i = (i + 1) % base.slotCount {
		//пытаемся положить в нужный слот задачу
		for index, val := range base.slots[i] {
			//как только натыкаемся на пустой элемент, кладем туда данные
			if len(val) == 0 {
				base.slots[i][index] = valueEvent
				//говорим что положили данные
				flag = true
				break
			}
		}
		//если данные положилы, прерываем попытки
		if flag {
			break
		}
	}
	if !flag {
		return errors.New("All slot is full")
	}
	return nil
}

func (base *TimerWheel) pinger() {
	for {
		base.ping <- true
		base.current = (base.current + 1) % base.slotCount
		time.Sleep(base.slotInterval)
	}
}

func (base *TimerWheel) handler() {
	for {
		<-base.ping
		if base.CallBack == nil {
			continue
		}
		for index, val := range base.slots[base.current] {
			if len(val) == 0 {
				break
			}
			base.CallBack(val)
			base.slots[base.current][index] = ""
		}
	}
}
