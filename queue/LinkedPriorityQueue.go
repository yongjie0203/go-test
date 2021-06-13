package main

import (

	"fmt"
	"math/rand"
	"time"
)

type IPriority interface {
	getPriority() (x int64)
}

type PriorityQueue struct{
	size int
	
	head Node
}

type Node struct{
	//last *Node

	data IPriority

	next *Node

}

func (this *PriorityQueue) insert(data IPriority ){
	var newNode = Node{}
	if(this.size == 0 ){
		newNode.data =  data
		newNode.next = nil
		this.head = newNode
		this.size++
		return 
	}
	newNode.data = data
	if(newNode.data.getPriority() > this.head.data.getPriority()){
		//this.head.last = &newNode
		this.head.next = nil
		newNode.next = &this.head
		this.head = newNode
		this.size ++

		fmt.Printf("queue is %v \n ",this)
		
	}else{
		var curNode = this.head
		for{
			fmt.Printf("curNode is : %v \n",curNode)
			if( nil != curNode.next && curNode.next.data.getPriority() > data.getPriority()){
				curNode = *curNode.next

			}else{
				
				break
			}
			
		}
		newNode.next = curNode.next
		curNode.next = &newNode
		this.size ++

	}

	
	
}




func (this *PriorityQueue) peek() (Node){

	return this.head

}


func (this *PriorityQueue) pop() (Node){
	node := Node{}
	if( this.size == 0){
		return node
	}
	var temp = this.head
	//this.head.next.last = nil
	this.head = *this.head.next
	this.size--
	return temp

}

type Order struct {
	Price int64 `json:"price"`
	Num   int   `json:"num"`
	Time  int64 `json:"time"`
}

func (this Order) getPriority() int64 {
	return this.Price * (time.Hour.Nanoseconds() - this.Time)
}


func main(){

	var queue = PriorityQueue{}

	fmt.Println("queue已初始化...")
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {

		var item Order
		item.Price = rand.Int63n(int64(100))
		item.Time = int64(i)
		item.Num = 1
		queue.insert(item)

		fmt.Printf("queue中插入%v \n", item)
		fmt.Printf("queue.size is %d \n", queue.size)
		//j, _ := json.Marshal(&item)
		//fmt.Printf("堆中插入jsonitem: %v\n", string(j))
		
		//fmt.Printf("%T堆中数据:%v  \n", maxHeap,maxHeap)

	}

	var curNode = queue.head
	var i = 0
	for {
		if(curNode.next != nil ){
			fmt.Printf("item in queue : %v \n ", curNode.next)
			i++
			if i >5 {break}
		}else{
			break
		}
	}
}