package main

import (
	"fmt"
	"sync"
)

var message chan Message = make(chan Message)
var count int = 0
var wg sync.WaitGroup

type Message struct {
	address    string
	privateKey string
	id         int
}

func main() {
	//var start = time.Now().Nanosecond()
	fmt.Printf("hhhh%d", 1)

	for i := 1; i <= 10000000; i++ {
		go pushMessage(i)

	}

	for {
		select {

		case v := <-message:
			printMessage(v)
		default:
			//wg.Wait()
			//var end = time.Now().Nanosecond()
			//fmt.Printf("耗时:%d", end-start)
			continue
		}

	}
}

func printMessage(ch Message) {
	fmt.Printf("message is %s count is %d \n", ch.address, ch.id)
	
}

func pushMessage(i int) {
	wg.Add(1)
	var m = new(Message)
	m.address = fmt.Sprintf("0x%d", i)
	m.privateKey = fmt.Sprintf("0x%d", i)
	count++
	m.id = count
	message <- *m
	wg.Done()
}
