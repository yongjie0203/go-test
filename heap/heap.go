package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

/*
	parent(i) = floor((i - 1)/2)
	left(i)   = 2i + 1
	right(i)  = 2i + 2
*/

type IPriority interface {
	getPriority() (x float64)
}

type Heap struct {
	/*堆容器*/
	heap []IPriority

	/* 当前堆元素个数 */
	size int

	/* 堆容器容量 */
	capacity int

	/*mark the heap is a max value root heap or min value root heap*/
	flag float64
}

const ROOT_VALUE_MAX float64 = 1
const ROOT_VALUE_MIN float64 = -1

/*将新元素压入堆中*/
func (this *Heap) insert(item IPriority) {
	if this.size == this.capacity {
		newHeap := make([]IPriority, this.capacity*2)
		copy(newHeap, this.heap)
		this.heap = newHeap
		this.capacity = this.capacity * 2
	}
	this.heap[this.size] = item
	this.size++
	this.heapUp()
}

/*将优先级最高的元素取出*/
func (this *Heap) poll() (root IPriority, err error) {
	if this.size == 0 {
		var item IPriority
		return item, errors.New("there is no node in heap")
	}
	root = this.heap[0]
	this.heap[0] = this.heap[this.size-1]
	this.size--
	this.heapDown()
	return root, nil
}

/*返回优先级最高的元素*/
func (this *Heap) peek() (IPriority, error) {
	if this.size == 0 {
		var item IPriority
		return item, errors.New("there is no node in heap")
	}

	return this.heap[0], nil
}

func (this *Heap) updateHead(item IPriority) {
	this.update(0, item)
}

func (this *Heap) update(index int, item IPriority) {
	if this.size-1 < index {
		return
	}
	if this.heap[index].getPriority() != item.getPriority() {
		panic("cant update ,priority is not match")
	}
	this.heap[index] = item
}

/* 初始化,设置堆的容量 */
func (this *Heap) initHeap(cap int) {
	this.capacity = cap
	this.flag = ROOT_VALUE_MAX //default heap has max value root
	this.heap = make([]IPriority, cap)
	this.size = 0
}

/* 通过节点索引获取该节点左子节点索引 */
func (this *Heap) getLeftChildIndex(parentIndex int) int {
	return 2*parentIndex + 1
}

/* 通过节点索引获取该节点右子节点索引 */
func (this *Heap) getRightChildIndex(parentIndex int) int {
	return this.getLeftChildIndex(parentIndex) + 1
}

/* 通过节点索引获取该节点父节点索引 */
func (this *Heap) getParentIndexByChildIndex(childIndex int) int {
	return (childIndex - 1) / 2
}

/* 是否存在左子节点 */
func (this *Heap) hasLeftChild(index int) bool {
	return this.getLeftChildIndex(index) < this.size
}

/* 是否存在右子节点 */
func (this *Heap) hasRightChild(index int) bool {
	return this.getRightChildIndex(index) < this.size
}

/* 是否存在父节点 */
func (this *Heap) hasParent(index int) bool {
	return this.getParentIndexByChildIndex(index) >= 0
}

/* 获取左子节点 */
func (this *Heap) leftChild(index int) IPriority {
	return this.heap[this.getLeftChildIndex(index)]
}

/* 获取右子节点 */
func (this *Heap) rightChild(index int) IPriority {
	return this.heap[this.getRightChildIndex(index)]
}

/* 获取父节点 */
func (this *Heap) parent(index int) IPriority {
	return this.heap[this.getParentIndexByChildIndex(index)]
}

/* 交换位置 */
func (this *Heap) swap(index1 int, index2 int) {
	this.heap[index1], this.heap[index2] = this.heap[index2], this.heap[index1]
}

func (this *Heap) isMaxRootHeap() bool {
	return this.flag > ROOT_VALUE_MIN
}

func (this *Heap) isMinRootHeap() bool {
	return this.flag == ROOT_VALUE_MIN
}

func (this *Heap) heapUp() {
	this.heapifyUp(this.size - 1)
}

func (this *Heap) heapDown() {
	this.heapifyDown(0)
}

func (this *Heap) heapifyUp(index int) {

	for {
		if this.hasParent(index) && this.parent(index).getPriority()*this.flag < this.heap[index].getPriority()*this.flag {

			this.swap(this.getParentIndexByChildIndex(index), index)

			index = this.getParentIndexByChildIndex(index)
		} else {
			break
		}

	}
}

func (this *Heap) heapifyDown(index int) {

	for {
		if this.hasLeftChild(index) {
			largerChindindex := this.getLeftChildIndex(index)
			if this.hasRightChild(index) && this.rightChild(index).getPriority()*this.flag > this.leftChild(index).getPriority()*this.flag {
				largerChindindex = this.getRightChildIndex(index)
			}

			if this.heap[index].getPriority()*this.flag < this.heap[largerChindindex].getPriority()*this.flag {
				this.swap(index, largerChindindex)
			} else {
				break
			}
			index = largerChindindex
		} else {
			break
		}
	}
}

func (this *Heap) delete(index int) {
	if this.size == 0 {
		return
	}
	//move last node to the node's index that has be delete
	this.heap[index] = this.heap[this.size-1]
	this.size--

	if this.isMaxRootHeap() {
		this.heapifyUp(index)
	} else {
		this.heapifyDown(index)
	}
}

/* 将堆拷贝到数组中 */
func (this Heap) copyAsArray() []IPriority {
	newArray := make([]IPriority, this.size)
	copy(newArray, this.heap[:this.size])
	return newArray
}

func (this Heap) copyAsSortArray() []IPriority {
	newHeap := Heap{}
	newHeap.heap = this.copyAsArray()
	newHeap.size = this.size
	newHeap.capacity = this.capacity
	newHeap.flag = this.flag
	array := make([]IPriority, newHeap.size)
	i := 0
	//fmt.Printf("newheap add is %x ,thisHeap addr is %x  \n" ,newHeap,this)
	for {
		if newHeap.size > 0 {
			array[i], _ = newHeap.poll()
			i++
		} else {
			break
		}
	}
	return array
}

func (this Heap) copyAsSortArrayLimit(index int) []IPriority {
	newHeap := Heap{}
	newHeap.heap = this.copyAsArray()
	newHeap.size = this.size
	newHeap.capacity = this.capacity
	newHeap.flag = this.flag
	array := make([]IPriority, index)
	i := 0
	//fmt.Printf("newheap add is %x ,thisHeap addr is %x  \n" ,newHeap,this)
	for {
		if newHeap.size > 0 && i < index {
			array[i], _ = newHeap.poll()
			i++
		} else {
			break
		}
	}
	return array
}

type Order struct {
	OID      int
	UID      int
	Price    float64 `json:"price"`
	Num      float64 `json:"num"`
	Time     int64   `json:"time"`
	Priority float64
}

func (this Order) getPriority() float64 {
	return this.Priority
}

func copyAsMapTopPriceLimit(arr []IPriority, limit int) map[float64]float64 {
	orderMap := make(map[float64]float64)
	keyCount := 0
	for i := 0; i < len(arr); i++ {

		switch t := arr[i].(type) {
		case Order:
			//fmt.Printf("type  %v \n", t.Price)
			num, ok := orderMap[t.Price]
			if ok {
				//fmt.Println(t.Price)
				orderMap[t.Price] = num + t.Num
			} else {

				orderMap[t.Price] = t.Num
				keyCount++

				if keyCount > limit {
					delete(orderMap, t.Price)
					break
				}
			}
		default:
			fmt.Println("unknown type")
		}

	}

	return orderMap
}

func peekOrder(heap *Heap) Order {
	var nilOrder Order
	if heap.size > 0 {
		top, _ := heap.peek()
		switch order := top.(type) {
		case Order:
			return order
		default:
			fmt.Println("unknown buyTopOrder type ")
			return nilOrder
		}
	} else {
		return nilOrder
	}

}

func pollOrder(heap *Heap) Order {
	var nilOrder Order
	if heap.size > 0 {
		top, _ := heap.poll()
		switch order := top.(type) {
		case Order:
			return order
		default:
			fmt.Println("unknown buyTopOrder type ")
			return nilOrder
		}
	} else {
		return nilOrder
	}

}

func onTransaction(sell, buy Order, transactionNum float64) {

	var content = `
	sellOid:%d	buyOid:%d
	sellUid:%d	buyUid:%d
	sellPrice:%f	buyPrice:%f
	sellhold:%f	buyhold:%f
	sellNum:%f	buyNum:%f
	`
	func() {
		log.Printf(content,
			sell.OID, buy.OID,
			sell.UID, buy.UID,
			sell.Price, buy.Price,
			sell.Num-transactionNum, buy.Num-transactionNum,
			transactionNum, transactionNum)

	}()

}

func transaction(sell, buy *Heap) {
	var count = 0

	for {
		if sell.size > 0 && buy.size > 0 {
			var sellTopOrder Order
			var buyTopOrder Order
			var transactionNum float64
			//var transactionPrice int

			sellTopOrder = peekOrder(sell)
			buyTopOrder = peekOrder(buy)

			if sellTopOrder.Price <= buyTopOrder.Price {
				//transactionPrice = int(sellTopOrder.Price)
				if sellTopOrder.Num <= buyTopOrder.Num {
					transactionNum = sellTopOrder.Num
				} else {
					transactionNum = buyTopOrder.Num
				}
				onTransaction(sellTopOrder, buyTopOrder, transactionNum)

				sellTopOrder.Num = sellTopOrder.Num - transactionNum
				buyTopOrder.Num = buyTopOrder.Num - transactionNum

				sell.updateHead(sellTopOrder)
				buy.updateHead(buyTopOrder)

				if sellTopOrder.Num <= 0 {
					//sell.poll()
					pollSell := pollOrder(sell)
					//test code
					if pollSell.OID != sellTopOrder.OID {
						log.Printf("sell Oid not match %v %v \n", pollSell.OID, sellTopOrder.OID)
						break
					}
					if pollSell.Num != sellTopOrder.Num {
						log.Printf("sell num not match %v %v \n", pollSell.Num, sellTopOrder.Num)
						break
					}
				}

				if buyTopOrder.Num <= 0 {
					//buy.poll()
					pollBuy := pollOrder(buy)
					//test code
					if pollBuy.OID != buyTopOrder.OID {
						log.Printf("buy Oid not match %v %v \n", pollBuy.OID, buyTopOrder.OID)
						break
					}
					if pollBuy.Num != buyTopOrder.Num {
						log.Printf("buy Num not match %v %v \n", pollBuy.Num, buyTopOrder.Num)
						break
					}
				}
				count++

				fmt.Printf("buy heap sise is %d sell heap size is %d\n", buy.size, sell.size)
				m := copyAsMapTopPriceLimit(buy.heap, 5)
				fmt.Printf("%s \n", formattedJson(&m))
				m = copyAsMapTopPriceLimit(sell.heap, 5)
				fmt.Printf("%s \n", formattedJson(&m))
			} else {
				time.Sleep(time.Nanosecond * 10000)
			}
		} else {
			time.Sleep(time.Nanosecond * 10000)
			//runtime.Gosched()
			//select {}
			//fmt.Print("all order is deal closed \n")
		}
	}
}

func buyTask(buyHeap *Heap) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {

		var item Order
		item.Price = rand.Float64() * float64(10)
		item.Time = time.Now().UnixNano()
		item.Num = rand.Float64() * float64(5)
		item.OID = i
		item.UID = rand.Intn(5)
		item.Priority = item.Price * float64(time.Now().UnixNano()-item.Time)
		buyHeap.insert(item)
		//fmt.Printf("buy Task: %v \n", item)
		//fmt.Printf("buy list : %v \n", copyAsMapTopPriceLimit(buyHeap.copyAsArray(), 5))
		//time.Sleep(time.Second * 5)

	}
}

func sellTask(sellHeap *Heap) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {

		var item Order
		item.Price = rand.Float64() * float64(10)
		item.Time = time.Now().UnixNano()
		item.Num = rand.Float64() * float64(5)
		item.OID = -1 * i
		item.UID = rand.Intn(5)
		item.Priority = item.Price * float64(time.Now().UnixNano()-item.Time)
		sellHeap.insert(item)
		//fmt.Printf("sell Task: %v \n", item)
		//fmt.Printf("sell list : %v \n", copyAsMapTopPriceLimit(sellHeap.copyAsArray(), 5))
		//time.Sleep(time.Second * 5)
	}
}

func formattedJson(obj *map[float64]float64) string {
	fmt.Printf("obj is %v :", obj)
	bs, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("bs is :" + string(bs))
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")
	fmt.Print("out string :" + out.String())
	return out.String()
}

func main() {

	var stopLock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		//阻塞程序运行，直到收到终止的信号
		var signal = <-signalChan
		stopLock.Lock()
		stop = true
		stopLock.Unlock()
		log.Printf("signal %v \n Cleaning before stop...", signal)
		stopChan <- struct{}{}
		os.Exit(0)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	buyHeap, sellHeap := Heap{}, Heap{}
	buyHeap.initHeap(5)
	sellHeap.initHeap(5)
	buyHeap.flag = ROOT_VALUE_MAX
	sellHeap.flag = ROOT_VALUE_MIN
	buyTask(&buyHeap)

	sellTask(&sellHeap)

	transaction(&sellHeap, &buyHeap)

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
