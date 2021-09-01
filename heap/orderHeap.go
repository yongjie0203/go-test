package heap

import (
	"errors"
	"fmt"
	"sync"

	"github.com/yongjie0203/go-trade-core/order"
)

/*
	parent(i) = floor((i - 1)/2)
	left(i)   = 2i + 1
	right(i)  = 2i + 2
*/

type OrderHeap struct {
	/*堆容器*/
	heap []order.Order

	/* 当前堆元素个数 */
	size int

	/* 堆容器容量 */
	capacity int

	/*mark the heap is a max value root heap or min value root heap*/
	flag float64

	insert *sync.RWMutex
}

func (h *OrderHeap) MaxHeap() {
	h.flag = ROOT_VALUE_MAX
}

func (h *OrderHeap) MinHeap() {
	h.flag = ROOT_VALUE_MIN
}

func (h *OrderHeap) Size() int {
	return h.size
}

func (h *OrderHeap) Lock() {
	h.insert.Lock()
}

func (h *OrderHeap) Unlock() {
	h.insert.Unlock()
}

/*将新元素压入堆中*/
func (h *OrderHeap) Insert(item order.Order) {
	fmt.Printf("call Insert %f \n", item.GetPriority())
	h.insert.Lock()
	defer h.insert.Unlock()
	if h.size == h.capacity {
		newHeap := make([]order.Order, h.capacity*2)
		copy(newHeap, h.heap)
		h.heap = newHeap
		h.capacity = h.capacity * 2
	}
	h.heap[h.size] = item
	h.size++
	h.heapUp()
}

/*将优先级最高的元素取出*/
func (h *OrderHeap) Poll() (root order.Order, err error) {
	h.insert.Lock()
	defer h.insert.Unlock()

	if h.size == 0 {
		var item order.Order
		return item, errors.New("there is no node in heap")
	}
	root = h.heap[0]
	fmt.Printf("call Poll %f \n", root.GetPriority())
	h.heap[0] = h.heap[h.size-1]
	h.size--
	h.heapDown()
	return root, nil
}

/*返回优先级最高的元素*/
func (this *OrderHeap) Peek() (*order.Order, error) {
	fmt.Println("call Peek")
	if this.size == 0 {
		var item order.Order
		return &item, errors.New("there is no node in heap")
	}

	return &this.heap[0], nil
}

func (this *OrderHeap) Head(item order.Order, callback func()) {
	//this.Lock()
	//defer this.Unlock()
	fmt.Printf("call Head %f \n", item.GetPriority())
	this.update(0, item, callback)
}

func (this *OrderHeap) update(index int, item order.Order, callback func()) {

	if this.size-1 < index {
		return
	}
	if this.heap[index].GetPriority() != item.GetPriority() {
		fmt.Printf("heap[%d].priority : %f  item.priority : %f \n", index, this.heap[index].GetPriority(), item.GetPriority())
		callback()
		//panic("cant update ,priority is not match")
	}
	this.heap[index] = item
}

/* 初始化,设置堆的容量 */
func (h *OrderHeap) InitHeap(cap int) {
	h.capacity = cap
	h.flag = ROOT_VALUE_MAX //default heap has max value root
	h.heap = make([]order.Order, cap)
	h.size = 0
	h.insert = new(sync.RWMutex)

}

/* 通过节点索引获取该节点左子节点索引 */
func (this *OrderHeap) getLeftChildIndex(parentIndex int) int {
	return 2*parentIndex + 1
}

/* 通过节点索引获取该节点右子节点索引 */
func (this *OrderHeap) getRightChildIndex(parentIndex int) int {
	return this.getLeftChildIndex(parentIndex) + 1
}

/* 通过节点索引获取该节点父节点索引 */
func (this *OrderHeap) getParentIndexByChildIndex(childIndex int) int {
	return (childIndex - 1) / 2
}

/* 是否存在左子节点 */
func (this *OrderHeap) hasLeftChild(index int) bool {
	return this.getLeftChildIndex(index) < this.size
}

/* 是否存在右子节点 */
func (this *OrderHeap) hasRightChild(index int) bool {
	return this.getRightChildIndex(index) < this.size
}

/* 是否存在父节点 */
func (this *OrderHeap) hasParent(index int) bool {
	return this.getParentIndexByChildIndex(index) >= 0
}

/* 获取左子节点 */
func (this *OrderHeap) leftChild(index int) order.Order {
	return this.heap[this.getLeftChildIndex(index)]
}

/* 获取右子节点 */
func (this *OrderHeap) rightChild(index int) order.Order {
	return this.heap[this.getRightChildIndex(index)]
}

/* 获取父节点 */
func (this *OrderHeap) parent(index int) order.Order {
	return this.heap[this.getParentIndexByChildIndex(index)]
}

/* 交换位置 */
func (this *OrderHeap) swap(index1 int, index2 int) {

	/*if index1 == 0 || index2 == 0 {

		func() {
			this.head.Lock()
			defer this.head.Unlock()
			this.heap[index1], this.heap[index2] = this.heap[index2], this.heap[index1]
			return
		}()

	}*/

	this.heap[index1], this.heap[index2] = this.heap[index2], this.heap[index1]
}

func (this *OrderHeap) isMaxRootHeap() bool {
	return this.flag > ROOT_VALUE_MIN
}

func (this *OrderHeap) isMinRootHeap() bool {
	return this.flag == ROOT_VALUE_MIN
}

func (this *OrderHeap) heapUp() {
	this.heapifyUp(this.size - 1)
}

func (this *OrderHeap) heapDown() {
	this.heapifyDown(0)
}

func (this *OrderHeap) heapifyUp(index int) {

	for {
		if this.hasParent(index) && this.parent(index).GetPriority()*this.flag < this.heap[index].GetPriority()*this.flag {

			this.swap(this.getParentIndexByChildIndex(index), index)

			index = this.getParentIndexByChildIndex(index)
		} else {
			break
		}

	}
}

func (this *OrderHeap) heapifyDown(index int) {

	for {
		if this.hasLeftChild(index) {
			largerChindindex := this.getLeftChildIndex(index)
			if this.hasRightChild(index) && this.rightChild(index).GetPriority()*this.flag > this.leftChild(index).GetPriority()*this.flag {
				largerChindindex = this.getRightChildIndex(index)
			}

			if this.heap[index].GetPriority()*this.flag < this.heap[largerChindindex].GetPriority()*this.flag {
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

func (this *OrderHeap) delete(index int) {
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
func (this OrderHeap) copyAsArray() []order.Order {
	newArray := make([]order.Order, this.size)
	copy(newArray, this.heap[:this.size])
	return newArray
}

func (this OrderHeap) copyAsSortArray() []order.Order {
	newHeap := OrderHeap{}
	newHeap.heap = this.copyAsArray()
	newHeap.size = this.size
	newHeap.capacity = this.capacity
	newHeap.flag = this.flag
	array := make([]order.Order, newHeap.size)
	i := 0
	//fmt.Printf("newheap add is %x ,thisHeap addr is %x  \n" ,newHeap,this)
	for {
		if newHeap.size > 0 {
			array[i], _ = newHeap.Poll()
			i++
		} else {
			break
		}
	}
	return array
}

func (this OrderHeap) copyAsSortArrayLimit(index int) []order.Order {
	newHeap := OrderHeap{}
	newHeap.heap = this.copyAsArray()
	newHeap.size = this.size
	newHeap.capacity = this.capacity
	newHeap.flag = this.flag
	array := make([]order.Order, index)
	i := 0
	//fmt.Printf("newheap add is %x ,thisHeap addr is %x  \n" ,newHeap,this)
	for {
		if newHeap.size > 0 && i < index {
			array[i], _ = newHeap.Poll()
			i++
		} else {
			break
		}
	}
	return array
}
