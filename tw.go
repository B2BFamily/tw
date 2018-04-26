package tw

import (
	"errors"
	"fmt"
	"time"
)

type TimerWheel struct {
	SlotInterval time.Duration      //Интервал времени обработки в секундах
	Duration     time.Duration      //Максимальное кол-во времени, которое будет обрабатываться
	MaxSlotEvent int                //максимальное кол-во событий на одну единицу времени
	CallBack     func(string) error //Обработчик событий
	slots        [][]string         //Колесо, в котором будет все храниться
	ping         chan bool          //канал
	current      int                //текущий номер элемента
	slotCount    int                //количество слотов в колесе
}

func (base *TimerWheel) Init() {
	base.current = 0
	if base.MaxSlotEvent == 0 {
		base.MaxSlotEvent = 20
	}
	base.ping = make(chan bool)
	base.slotCount = int(base.Duration.Seconds() / base.SlotInterval.Seconds())
	base.slots = make([][]string, base.slotCount)
	for i := 0; i < base.slotCount; i++ {
		base.slots[i] = make([]string, base.MaxSlotEvent)
	}
	go base.pinger()
	go base.handler()
}

func (base *TimerWheel) Set(timeEvent time.Time, valueEvent string) error {
	var duration = time.Duration(timeEvent.UTC().Unix()-time.Now().UTC().Unix()) * time.Second
	slotIndex := int(duration.Seconds() / base.SlotInterval.Seconds())
	if slotIndex > base.slotCount {
		return errors.New(fmt.Sprintf("time of event is too long, event time: %v, event duration: %v, max duration: %v", timeEvent, duration, base.Duration))
	}
	slotIndex = (slotIndex + base.current) % base.slotCount
	flag := false
	for i := slotIndex; i != slotIndex-1; i++ {
		i = i % base.slotCount
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
	return nil
}

//с указанным интервалом пингует обработчик и говорит, в каком элементе массива смотреть события
func (base *TimerWheel) pinger() {
	for {
		base.ping <- true
		base.current++
		if base.current >= base.slotCount {
			base.current = 0
		}
		time.Sleep(base.SlotInterval)
	}
}

//слушает пингер и отдает все элементы
func (base *TimerWheel) handler() {
	for {
		<-base.ping
		fmt.Println(base.current)
		fmt.Println(time.Now().UTC().Unix())
		if base.CallBack == nil {
			continue
		}
		for index, val := range base.slots[base.current] {
			if len(val) == 0 {
				break
			}
			fmt.Println(val)
			base.CallBack(val)
			base.slots[base.current][index] = ""
		}
	}
}
