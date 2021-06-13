package main

import (

	"fmt"
	"math/rand"
	"time"
)

type IPriority interface {
	getPriority() (x int64)
}

type ArrayPriorityQueue struct{

	index int
	
	data []*IPriority
}



func (this *ArrayPriorityQueue) insert(item IPriority ){
	var isAdd = false
	if(len(this.data) == 0){
		//this.data = make([]*IPriority,1,10)
		this.data = append(this.data, &item) 
		//fmt.Printf("1 item : %v  %v\n",&item,item)
		this.index = 0
		return
	}

	//fmt.Printf("queue data is %v \n " ,toObjectArray(this.data))
	
	for i:= 0 ;i<len(this.data) ; i++{
		if(this.data[i] == nil) {
			continue
		}
		if(item.getPriority() > IPriority(*this.data[i]).getPriority() ){

			tail := make([]*IPriority ,len(this.data) -i)
			head := make([]*IPriority, i)
			copy(tail,this.data[i:])
			copy(head,this.data[this.index:i + this.index ])

			//fmt.Printf("2 item : %v  %v\n",&item,item)
			
			

			this.data=  append(head,&item)

			//fmt.Printf("head append item is %v \n",toObjectArray(this.data))
			this.data = append(this.data,tail...)
			//fmt.Printf("data  append tail is %v \n",toObjectArray(this.data))
			
			this.index = 0
			isAdd = true
			break
		}
	}

	if( !isAdd){
		this.data =  append(this.data,&item)
		//fmt.Printf("3 item  : %v  %v\n",&item,item)
	}
	
	//fmt.Printf("data addr %x  cap: %d len: %d\n",this.data,cap(this.data),len(this.data))
	
}




func (this *ArrayPriorityQueue) peek() ( IPriority){

	return *this.data[this.index]

}


func (this *ArrayPriorityQueue) pop() (IPriority){
	var item = this.data[this.index]
	this.index ++ 
	return *item
}

type Order struct {
	Price int64 `json:"price"`
	Num   int   `json:"num"`
	Time  int64 `json:"time"`
}

func (this Order) getPriority() int64 {
	return this.Price //* (time.Hour.Nanoseconds() - this.Time)
}

func toObjectArray(arr []*IPriority) []IPriority{
	list := make([]IPriority,len(arr))
	for i:= 0;i< len(arr);i++{
		if (arr[i] == nil){
			list[i] = nil
			continue
		}
		list[i] = *arr[i]		
	}
	return list 
}


func main(){

	var queue = ArrayPriorityQueue{}

	fmt.Println("queue已初始化...")
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000000; i++ {

		var item Order
		item.Price = rand.Int63n(int64(10))
		item.Time = int64(i)
		item.Num = 0
		queue.insert(item)

		if i % 10000 == 0{
			fmt.Printf("P: %d % \n",i/10000)
		}

		//fmt.Printf("queue中插入%v \n", item)
		//fmt.Printf("queue.size is %d \n", queue.size)
		//j, _ := json.Marshal(&item)
		//fmt.Printf("堆中插入jsonitem: %v\n", string(j))
		
		//fmt.Printf("%T堆中数据:%v  \n", maxHeap,maxHeap)

	}

	fmt.Printf("queue中插入 \n")

	//fmt.Printf("over queue data is %v \n" ,toObjectArray(queue.data))
	

	for{
		fmt.Printf("queue size : %d   index  %d \n",len(queue.data),queue.index)
		if queue.index < len(queue.data){
			queue.pop()
			//fmt.Printf("poll item %v ",item)
		}else{
			break
		}
		
		
	}

	fmt.Printf("queue中pop \n")

	for i := 0; i < 5; i++ {

		var item Order
		item.Price = rand.Int63n(int64(10))
		item.Time = int64(i)
		item.Num = 0
		queue.insert(item)

		//fmt.Printf("queue中插入%v \n", item)
		//fmt.Printf("queue.size is %d \n", queue.size)
		//j, _ := json.Marshal(&item)
		//fmt.Printf("堆中插入jsonitem: %v\n", string(j))
		
		//fmt.Printf("%T堆中数据:%v  \n", maxHeap,maxHeap)

	}
	//fmt.Printf("over queue data is %v \n" ,toObjectArray(queue.data))


}