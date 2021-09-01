package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"yingyi.cn/go-trade-core/heap"
	"yingyi.cn/go-trade-core/order"
	"yingyi.cn/go-trade-core/trade"
)

func buyTask(buyHeap *heap.Heap) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {

		var item order.Order
		item.Price = float64(i) //rand.Float64() * float64(10)
		item.Time = time.Now().UnixNano()
		item.Num = rand.Float64() * float64(5)
		item.OID = int64(i)
		item.UID = rand.Int63n(5)
		item.Priority = item.Price //+ float64(time.Now().UnixNano()-item.Time)
		buyHeap.Insert(item)
		//fmt.Printf("buy Task: %v \n", item)
		//fmt.Printf("buy list : %v \n", copyAsMapTopPriceLimit(buyHeap.copyAsArray(), 5))
		//time.Sleep(time.Second * 5)

	}
}

func sellTask(sellHeap *heap.Heap) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {

		var item order.Order
		item.Price = float64(i) //rand.Float64() * float64(10)
		item.Time = time.Now().UnixNano()
		item.Num = rand.Float64() * float64(3)
		item.OID = int64(-1 * i)
		item.UID = rand.Int63n(5)
		item.Priority = item.Price //+ float64(time.Now().UnixNano()-item.Time)
		sellHeap.Insert(item)
		//fmt.Printf("sell Task: %v \n", item)
		//fmt.Printf("sell list : %v \n", copyAsMapTopPriceLimit(sellHeap.copyAsArray(), 5))
		//time.Sleep(time.Second * 5)
	}
}

func main() {
	var start = time.Now().UnixNano()
	buyHeap, sellHeap := heap.Heap{}, heap.Heap{}

	var stopLock sync.Mutex
	//stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		//阻塞程序运行，直到收到终止的信号
		var signal = <-signalChan
		stopLock.Lock()
		//stop = true
		stopLock.Unlock()
		log.Printf("signal %v \n Cleaning before stop...", signal)
		b := buyHeap.CopyAsMapTopPriceLimit(5)
		s := sellHeap.CopyAsMapTopPriceLimit(5)
		sellHeap.SortByPrice(s)
		fmt.Println("------------------------------------------------------")
		buyHeap.SortByPrice(b)

		//fmt.Printf("\n buy list %s \n sell list %s \n", heap.FormattedJson(&b), heap.FormattedJson(&s))

		fmt.Printf("sellHeap.size %d buyHeap.size %d \n", sellHeap.Size(), buyHeap.Size())
		stopChan <- struct{}{}
		var end = time.Now().UnixNano()
		var tins = float64(trade.Count) / float64(end-start) * float64(1000000)
		fmt.Printf("profermance:%f \n", tins)
		os.Exit(0)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	buyHeap.InitHeap(5)
	sellHeap.InitHeap(5)
	sellHeap.MinHeap()
	buyHeap.MaxHeap()

	buyTask(&buyHeap)

	sellTask(&sellHeap)

	for i := 0; i < 1; i++ {

		go trade.Transaction(&sellHeap, &buyHeap)
	}

	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for _ = range ticker.C {
			b := buyHeap.CopyAsMapTopPriceLimit(5)
			s := sellHeap.CopyAsMapTopPriceLimit(5)
			fmt.Println("*****************************************************")
			sellHeap.SortByPrice(s)
			fmt.Println("------------------------------------------------------")
			buyHeap.SortByPrice(b)
			fmt.Println("*****************************************************")

			fmt.Printf("sellHeap.size %d buyHeap.size %d \n", sellHeap.Size(), buyHeap.Size())
		}
	}()

	for i := 0; i < 1000; i++ {

		buyTask(&buyHeap)

		sellTask(&sellHeap)
	}

	//panic("data is added")
	time.Sleep(time.Minute * 10)

	/*maxHeap, minHeap := Heap{}, Heap{}
	maxHeap.initHeap(5)
	minHeap.initHeap(5)
	maxHeap.flag = ROOT_VALUE_MAX
	minHeap.flag = ROOT_VALUE_MIN
	fmt.Println("堆已初始化...")
	rand.Seed(time.Now().UnixNano())
	var start = time.Now().UnixNano()

	for i := 0; i < 1000000; i++ {

		var item Order
		item.Price = rand.Int63n(int64(100))
		item.Time = int64(i)
		item.Num = 1
		maxHeap.insert(item)
		//j, _ := json.Marshal(&item)
		//fmt.Printf("堆中插入jsonitem: %v\n", string(j))
		//fmt.Printf("堆中插入%v \n", item)
		//fmt.Printf("%T堆中数据:%v  \n", maxHeap,maxHeap)

	}
	var end = time.Now().UnixNano()
	fmt.Printf("time1 : %v \n", (end-start)/1e6)

	maxHeap.copyAsSortArray()
	//fmt.Printf("堆拷贝:%v \n", maxHeap.copyAsArray())
	//fmt.Printf("sort堆拷贝:%v \n", maxHeap.copyAsSortArray())

	var end2 = time.Now().UnixNano()

	var item Order
	item.Price = rand.Int63n(int64(100))
	item.Time = int64(100000)
	item.Num = 1
	maxHeap.insert(item)

	var end3 = time.Now().UnixNano()

	fmt.Printf("time2 : %v \n", (end2-end)/1e6)
	fmt.Printf("time3 : %v \n", (end3-end2)/1e6)

	fmt.Printf("limit soft data : %v", copyAsMapTopPriceLimit(maxHeap.copyAsArray(), 5))
	//fmt.Printf("堆拷贝:%v \n", maxHeap.copyAsArray())
	*/
	/*for i := 0; i < 6; i++ {

		var item HeapItem
		item.priority = rand.Intn(100)
		minHeap.insert(item)
		fmt.Printf("堆中插入%v \n", item)
		//fmt.Printf("堆中数据:%v \n", minHeap)
		fmt.Printf("堆拷贝:%v \n", minHeap.copyAsArray())

	}*/

	/*for i := 0; i < 6; i++ {

		fmt.Printf("弹出堆顶数据%d \n", poll())
		fmt.Printf("堆中数据:%v \n", heap)

		fmt.Printf("堆拷贝:%v \n", maxHeap.copayAsSortArray())
	}*/

	//delete(2)

}
