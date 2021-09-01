package queue

import (
	"fmt"
	"log"
	//"math/rand"
	//"time"
	"bytes"
	"encoding/json"
	"sync"
	//"runtime"
	//"os"
	"github.com/yongjie0203/go-trade-core/order"
)

type PriorityQueue struct {
	size int

	head *Node

	tail *Node

	lock *sync.RWMutex

	tailLock *sync.RWMutex

	headLock *sync.RWMutex

	name string
}

type Node struct {
	prev *Node

	data order.IPriority

	next *Node

	lock *sync.RWMutex
}

func (n *Node) UpdateData(data order.IPriority) {
	n.data = data
}

func (q *PriorityQueue) Remove(node *Node) {
	node.lock.Lock()
	defer node.lock.Unlock()
	if node.prev == nil {
		if node.next != nil {
			node.next.prev = nil
		}
		q.head = node.next
		if q.head != nil {
			q.head.prev = nil
		}
		q.size--
		return
	}
	if node.next == nil {
		if node.prev != nil {
			node.prev.next = nil
		}
		q.tail = node.prev
		if q.tail != nil {
			q.tail.next = nil
		}
		q.size--
		return
	}

	node.prev.next = node.next
	node.next.prev = node.prev
	q.size--
	return
}

var h = 0

func (q *PriorityQueue) Insert(data order.IPriority) *Node {
	q.lock.Lock()
	defer q.lock.Unlock()
	isInsert := false
	if q.size == 0 {
		q.UpdateHead(data)
		q.size++
		isInsert = true
		//PrintQueue(*q)
		return q.head
	}

	var node = &Node{data: data, next: nil}
	node.lock = new(sync.RWMutex)

	for cur := q.head; cur != nil; cur = q.succ(cur) {

		if data.GetPriority() > cur.data.GetPriority() {
			//fmt.Printf("data: %v  cur:%v \n",data.GetPriority(),cur.data.GetPriority())

			node.prev = cur.prev
			cur.prev = node
			node.next = cur

			if node.prev != nil {
				node.prev.next = node

			} else {

				q.head = node

			}
			q.size++
			isInsert = true

			//cur = node
			//fmt.Printf("for [%v]\n",cur.data.GetPriority())
			//PrintQueue(*q)
			return node
		}

	}

	if !isInsert {

		node.prev = q.tail
		node.next = nil
		q.tail.next = node
		q.tail = node
		q.size++
		isInsert = true
		//node.prev.lock.Unlock()
		//fmt.Printf("tail [%v]\n",data.GetPriority())
		//PrintQueue(*q)
		return q.tail
	}
	PrintQueue(*q)
	return node
}

func (q *PriorityQueue) succ(p *Node) *Node {
	var next = p.next
	if p == next {
		return q.head
	} else {
		return next
	}
}

func (q *PriorityQueue) UpdateHead(data order.IPriority) {
	q.headLock.Lock()
	defer q.headLock.Unlock()
	if q.size == 0 {
		var node = &Node{data: data, next: nil}
		node.lock = new(sync.RWMutex)

		q.head = node

		q.tail = node

	} else {

		var node = &Node{data: data}
		node.lock = new(sync.RWMutex)

		node.prev = nil
		node.next = q.head.next
		q.head.next.prev = node

		q.head = node
	}
}

func (this *PriorityQueue) Peek() *Node {

	return this.head

}

func (q *PriorityQueue) Pop() *Node {
	q.lock.Lock()
	defer q.lock.Unlock()
	node := &Node{}
	if q.size == 0 {
		return node
	}
	var temp = q.head
	//this.head.next.last = nil
	q.head = q.head.next
	q.head.prev = nil
	q.size--
	return temp

}

func (q *PriorityQueue) Size() int {
	size := 0
	for cur := q.head; cur != nil; cur = q.succ(cur) {
		size++
	}
	return size
}

func CheckQueue(q PriorityQueue) int {
	size := 0
	for cur := q.head; cur != nil; cur = q.succ(cur) {
		var prev = cur.prev
		if prev != nil {
			if prev.data.GetPriority() < cur.data.GetPriority() {
				fmt.Printf("queue:%s  prev:%v,cur:%v\n", q.name, prev.data.GetPriority() < cur.data.GetPriority())
			}
		}
		size++
	}
	return size
}

func PrintQueue(queue PriorityQueue) {

	fmt.Printf("\n %s \n", queue.name)
	for cur := queue.head; cur != nil; cur = queue.succ(cur) {
		fmt.Printf("[%v]->", cur.data.GetPriority())
	}
	fmt.Printf("\n")
}

func (q PriorityQueue) CopyAsSortArrayLimit(limit int) {
	var arr = make([]order.IPriority, 100)
	i := 0
	var prevPrice float64 = 0
	for cur := q.head; cur != nil; cur = q.succ(cur) {
		switch order := cur.data.(type) {
		case order.Order:
			if order.Price != prevPrice {
				prevPrice = order.Price
				i++
			}
			if i <= limit {
				arr = append(arr, order)
			} else {
				break
			}

		default:

		}
	}

}

func (q *PriorityQueue) Name(name string) *PriorityQueue {
	q.name = name
	return q
}

func (q *PriorityQueue) SetLock(lock *sync.RWMutex) *PriorityQueue {
	q.lock = lock
	return q
}

func (q *PriorityQueue) SetHeadLock(lock *sync.RWMutex) *PriorityQueue {
	q.headLock = lock
	return q
}

func (q *PriorityQueue) SetTailLock(lock *sync.RWMutex) *PriorityQueue {
	q.tailLock = lock
	return q
}

func (q PriorityQueue) CopyAsMapTopPriceLimit(limit int) map[string]float64 {
	orderMap := make(map[string]float64)
	keyCount := 0
	for cur := q.head; cur != nil; cur = q.succ(cur) {

		switch t := cur.data.(type) {
		case order.Order:
			//fmt.Printf("type  %v \n", t.Price)
			num, ok := orderMap[fmt.Sprintf("%v", t.Price)]
			if ok {
				//fmt.Println(t.Price)
				orderMap[fmt.Sprintf("%v", t.Price)] = num + t.Num
			} else {

				orderMap[fmt.Sprintf("%v", t.Price)] = t.Num
				keyCount++

				if keyCount > limit {
					delete(orderMap, fmt.Sprintf("%v", t.Price))
					break
				}
			}
		default:
			fmt.Printf("unknown type %T \n", t)
		}

	}

	return orderMap
}

func FormattedJson(obj interface{}) string {
	//fmt.Printf("obj is %v :", obj)
	bs, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Print("bs is :" + string(bs))
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")
	//fmt.Print("out string :" + out.String())
	return out.String()
}

/**
func buyTask(buyQueue *PriorityQueue) {
	defer func(){
		//log.Println(recover())
	}()

	wg.Add(1)
	//rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1; i++ {

		ddlock.Lock()
		dd++
		var item Order
		item.Price = int64(dd)//rand.Int63n(int64(100))//float64(i) //rand.Float64() * float64(10)
		ddlock.Unlock()
		item.Time = time.Now().UnixNano()
		item.Num = float64(i+1)*0.1//rand.Float64() * float64(5)
		item.OID = i
		item.UID = i//rand.Intn(5)
		item.Priority = item.Price //+ float64(time.Now().UnixNano()-item.Time)
		buyQueue.Insert(item)
		//fmt.Printf("buy Task: %v \n", item)
		//fmt.Printf("buy list : %v \n", copyAsMapTopPriceLimit(buyHeap.copyAsArray(), 5))
		//time.Sleep(time.Second * 5)

	}
	wg.Done()
}

func sellTask(sellQueue *PriorityQueue) {
	defer func(){
		//log.Println(recover())
	}()
	wg.Add(1)
	//rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1; i++ {

		ddlock.Lock()
		dd++
		var item Order
		item.Price = int64(dd)//rand.Int63n(int64(100))//float64(i) //rand.Float64() * float64(10)
		ddlock.Unlock()
		item.Time = time.Now().UnixNano()
		item.Num = float64(i+1)*0.1//rand.Float64() * float64(3)
		item.OID = -1 * i
		item.UID =  i//rand.Intn(5)
		item.Priority = item.Price *-1//+ float64(time.Now().UnixNano()-item.Time)
		sellQueue.Insert(item)
		//fmt.Printf("sell Task: %v \n", item)
		//fmt.Printf("sell list : %v \n", copyAsMapTopPriceLimit(sellHeap.copyAsArray(), 5))
		//time.Sleep(time.Second * 5)
	}
	wg.Done()
}

var wg sync.WaitGroup
var dd = 0
var ddlock *sync.RWMutex

func main(){

	logfile,err := os.OpenFile("log.txt",os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer logfile.Close()
	defer func(){
		//log.Println(recover())
	}()
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}

	log.SetOutput(logfile)//设置输出位置

	var buyQueue = PriorityQueue{}
	var sellQueue = PriorityQueue{}
	buyQueue.name = "buy"
	sellQueue.name = "sell"
	buyQueue.lock = new(sync.RWMutex)
	sellQueue.lock = new(sync.RWMutex)
	buyQueue.tailLock = new(sync.RWMutex)
	sellQueue.tailLock = new(sync.RWMutex)
	buyQueue.headLock = new(sync.RWMutex)
	sellQueue.headLock = new(sync.RWMutex)

	ddlock =  new(sync.RWMutex)
	//buyTask(&buyQueue)

	//sellTask(&sellQueue)

	for i := 0; i < 2; i++ {

		// Transaction(&sellQueue, &buyQueue)
	}
	now := time.Now()
	for i := 0; i < 100000; i++ {

	go	buyTask(&buyQueue)

	go	sellTask(&sellQueue)
	}

	//panic("data is added")
	time.Sleep(time.Second * 3)
	wg.Wait()
	//println(runtime.NumGoroutine())
	fmt.Printf("\n")
	fmt.Print(buyQueue.size)
	fmt.Printf("\n")
	//PrintQueue(buyQueue)
	//fmt.Printf("buy:size : %d\n",buyQueue.Size())
	//PrintQueue(sellQueue)
	fmt.Print(sellQueue.size)
	//fmt.Printf("sell:size : %d\n",sellQueue.Size())
	duration := time.Since(now)
	fmt.Println("duration:", duration)
	now = time.Now()
	fmt.Printf("buy:size : %d\n",CheckQueue(buyQueue))
	duration = time.Since(now)
	fmt.Println("duration:", duration)

	//var node *Node
	/*
	fmt.Println("queue已初始化...")
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {

		var item Order
		item.Price = rand.Int63n(int64(100))
		item.Time = int64(i)
		item.Num = 1
		curNode := queue.Insert(item)

		if(i==0){
			node = curNode
		}

		fmt.Printf("queue中插入%v \n", item)
		fmt.Printf("queue.size is %d \n", queue.size)
		//j, _ := json.Marshal(&item)
		//fmt.Printf("堆中插入jsonitem: %v\n", string(j))

		//fmt.Printf("%T堆中数据:%v  \n", maxHeap,maxHeap)

	}

	PrintQueue(queue)
	queue.Remove(node)
	PrintQueue(queue)

	queue.Pop()
	queue.Pop()
	PrintQueue(queue)
*/

//}

func peekOrder(queue *PriorityQueue) (order.Order, *Node) {
	var nilOrder order.Order
	if queue.size > 0 {
		top := queue.Peek()
		switch order := top.data.(type) {
		case order.Order:
			//fmt.Printf("peek addr : %x  order addr : %x \n", &top, order)
			return order, top
		default:
			fmt.Printf("unknown peekOrder type %T %v \n", order, order)
			return nilOrder, top
		}
	} else {
		return nilOrder, nil
	}

}

func Transaction(sell, buy *PriorityQueue) {

	for {
		if buy.size > 0 && sell.size > 0 {
			buyTopOrder, buyNode := peekOrder(buy)
			sellTopOrder, sellNode := peekOrder(sell)
			if buyTopOrder.Price >= sellTopOrder.Price {
				var transactionNum float64
				if sellTopOrder.Num <= buyTopOrder.Num {
					transactionNum = sellTopOrder.Num
				} else {
					transactionNum = buyTopOrder.Num
				}

				onTransaction(sellTopOrder, buyTopOrder, transactionNum)
				sellTopOrder.Num = sellTopOrder.Num - transactionNum
				buyTopOrder.Num = buyTopOrder.Num - transactionNum
				buyNode.data = buyTopOrder
				sellNode.data = sellTopOrder

				if sellTopOrder.Num <= 0 {
					sell.Remove(sellNode)
				}
				if buyTopOrder.Num <= 0 {
					buy.Remove(buyNode)
				}

				fmt.Printf("sellQueue:\n")
				fmt.Printf("%v\n", FormattedJson(sell.CopyAsMapTopPriceLimit(5)))
				fmt.Printf("buyQueue:\n")
				fmt.Printf("%v\n", FormattedJson(buy.CopyAsMapTopPriceLimit(5)))

			}
			//no order to deal

		} else {
			//queue is empty
		}

	}
}

type HandlerTransaction func(sell, buy order.Order, transactionNum float64)

var registeredTransactionHandler = make(map[string]HandlerTransaction)

// RegistTransactionHandler 注册的处理器将并行执行
func RegistTransactionHandler(name string, handler HandlerTransaction) {
	registeredTransactionHandler[name] = handler
}

func onTransaction(sell, buy order.Order, transactionNum float64) {

	for name, handler := range registeredTransactionHandler {
		log.Printf("call registed transaction handler %s", name)
		go handler(sell, buy, transactionNum)
	}

	//Count++
	var content = `
	sellOid:%d	buyOid:%d
	sellUid:%d	buyUid:%d
	sellPrice:%v	buyPrice:%v
	sellhold:%v	buyhold:%v
	sellNum:%v	buyNum:%v
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
