package heap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/yongjie0203/go-trade-core/order"
)

/*
	parent(i) = floor((i - 1)/2)
	left(i)   = 2i + 1
	right(i)  = 2i + 2
*/

type Heap struct {
	/*堆容器*/
	heap []IPriority

	/* 当前堆元素个数 */
	size int

	/* 堆容器容量 */
	capacity int

	/*mark the heap is a max value root heap or min value root heap*/
	flag float64

	insert *sync.RWMutex
}

func (h *Heap) MaxHeap() {
	h.flag = ROOT_VALUE_MAX
}

func (h *Heap) MinHeap() {
	h.flag = ROOT_VALUE_MIN
}

func (h *Heap) Size() int {
	return h.size
}

func (h *Heap) Lock() {
	h.insert.Lock()
}

func (h *Heap) Unlock() {
	h.insert.Unlock()
}

const ROOT_VALUE_MAX float64 = 1
const ROOT_VALUE_MIN float64 = -1

/*将新元素压入堆中*/
func (h *Heap) Insert(item IPriority) {
	//fmt.Printf("call Insert %f \n", item.GetPriority())
	h.insert.Lock()
	defer h.insert.Unlock()
	if h.size == h.capacity {
		newHeap := make([]IPriority, h.capacity*2)
		copy(newHeap, h.heap)
		h.heap = newHeap
		h.capacity = h.capacity * 2
	}
	h.heap[h.size] = item
	h.size++
	h.heapUp()
}

/*将优先级最高的元素取出*/
func (h *Heap) Poll() (root IPriority, err error) {
	//h.insert.Lock()
	//defer h.insert.Unlock()

	if h.size == 0 {
		var item IPriority
		return item, errors.New("there is no node in heap")
	}
	root = h.heap[0]
	//	fmt.Printf("call Poll %f \n", root.GetPriority())
	h.heap[0] = h.heap[h.size-1]
	h.size--
	h.heapDown()
	return root, nil
}

/*返回优先级最高的元素*/
func (h *Heap) Peek() (IPriority, error) {
	//fmt.Println("call Peek")
	if h.size == 0 {
		var item IPriority
		return item, errors.New("there is no node in heap")
	}

	return h.heap[0], nil
}

func (h *Heap) Head(item IPriority, callback func()) {
	//h.Lock()
	//defer h.Unlock()
	fmt.Printf("call Head %f \n", item.GetPriority())
	h.update(0, item, callback)
}

func (h *Heap) update(index int, item IPriority, callback func()) {

	if h.size-1 < index {
		return
	}
	if h.heap[index].GetPriority() != item.GetPriority() {
		fmt.Printf("heap[%d].priority : %f  item.priority : %f \n", index, h.heap[index].GetPriority(), item.GetPriority())
		//callback()
		panic("cant update ,priority is not match")
	}
	h.heap[index] = item
}

/* 初始化,设置堆的容量 */
func (h *Heap) InitHeap(cap int) {
	h.capacity = cap
	h.flag = ROOT_VALUE_MAX //default heap has max value root
	h.heap = make([]IPriority, cap)
	h.size = 0
	h.insert = new(sync.RWMutex)

}

/* 通过节点索引获取该节点左子节点索引 */
func (h *Heap) getLeftChildIndex(parentIndex int) int {
	return 2*parentIndex + 1
}

/* 通过节点索引获取该节点右子节点索引 */
func (h *Heap) getRightChildIndex(parentIndex int) int {
	return h.getLeftChildIndex(parentIndex) + 1
}

/* 通过节点索引获取该节点父节点索引 */
func (h *Heap) getParentIndexByChildIndex(childIndex int) int {
	return (childIndex - 1) / 2
}

/* 是否存在左子节点 */
func (h *Heap) hasLeftChild(index int) bool {
	return h.getLeftChildIndex(index) < h.size
}

/* 是否存在右子节点 */
func (h *Heap) hasRightChild(index int) bool {
	return h.getRightChildIndex(index) < h.size
}

func (h *Heap) hasChild(index int) bool {
	return h.hasLeftChild(index) || h.hasRightChild(index)
}

func (h *Heap) hasDoubleChild(index int) bool {
	return h.hasLeftChild(index) && h.hasRightChild(index)
}

func (h *Heap) Next(index int) (iPriority IPriority, i int) {
	var item IPriority
	if !h.hasChild(index) {
		return item, -1
	}

	if h.hasDoubleChild(index) {
		if h.leftChild(index).GetPriority()*h.flag >= h.rightChild(index).GetPriority()*h.flag {
			return h.leftChild(index), h.getLeftChildIndex(index)
		} else {
			return h.rightChild(index), h.getRightChildIndex(index)
		}
	}

	if h.hasLeftChild(index) {
		return h.leftChild(index), h.getLeftChildIndex(index)
	}

	if h.hasRightChild(index) {
		return h.rightChild(index), h.getRightChildIndex(index)
	}
	return item, -1
}

/* 是否存在父节点 */
func (h *Heap) hasParent(index int) bool {
	return h.getParentIndexByChildIndex(index) >= 0
}

/* 获取左子节点 */
func (h *Heap) leftChild(index int) IPriority {
	return h.heap[h.getLeftChildIndex(index)]
}

/* 获取右子节点 */
func (h *Heap) rightChild(index int) IPriority {
	return h.heap[h.getRightChildIndex(index)]
}

/* 获取父节点 */
func (h *Heap) parent(index int) IPriority {
	return h.heap[h.getParentIndexByChildIndex(index)]
}

/* 交换位置 */
func (h *Heap) swap(index1 int, index2 int) {

	h.heap[index1], h.heap[index2] = h.heap[index2], h.heap[index1]
}

func (h *Heap) isMaxRootHeap() bool {
	return h.flag > ROOT_VALUE_MIN
}

func (h *Heap) isMinRootHeap() bool {
	return h.flag == ROOT_VALUE_MIN
}

func (h *Heap) heapUp() {
	h.heapifyUp(h.size - 1)
}

func (h *Heap) heapDown() {
	h.heapifyDown(0)
}

func (h *Heap) heapifyUp(index int) {

	for {
		if h.hasParent(index) && h.parent(index).GetPriority()*h.flag < h.heap[index].GetPriority()*h.flag {

			h.swap(h.getParentIndexByChildIndex(index), index)

			index = h.getParentIndexByChildIndex(index)
		} else {
			break
		}

	}
}

func (h *Heap) heapifyDown(index int) {

	for {
		if h.hasLeftChild(index) {
			largerChindindex := h.getLeftChildIndex(index)
			if h.hasRightChild(index) && h.rightChild(index).GetPriority()*h.flag > h.leftChild(index).GetPriority()*h.flag {
				largerChindindex = h.getRightChildIndex(index)
			}

			if h.heap[index].GetPriority()*h.flag < h.heap[largerChindindex].GetPriority()*h.flag {
				h.swap(index, largerChindindex)
			} else {
				break
			}
			index = largerChindindex
		} else {
			break
		}
	}
}

func (h *Heap) delete(index int) {
	if h.size == 0 {
		return
	}
	//move last node to the node's index that has be delete
	h.heap[index] = h.heap[h.size-1]
	h.size--

	if h.isMaxRootHeap() {
		h.heapifyUp(index)
	} else {
		h.heapifyDown(index)
	}
}

/* 将堆拷贝到数组中 */
func (h Heap) CopyAsArray() []IPriority {
	newArray := make([]IPriority, h.size)
	copy(newArray, h.heap[:h.size])
	return newArray
}

func (h Heap) CopyAsSortArray() []IPriority {
	newHeap := Heap{}
	newHeap.heap = h.CopyAsArray()
	newHeap.size = h.size
	newHeap.capacity = h.capacity
	newHeap.flag = h.flag
	array := make([]IPriority, newHeap.size)
	i := 0
	//fmt.Printf("newheap add is %x ,hHeap addr is %x  \n" ,newHeap,h)
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

func (h Heap) CopyAsSortArrayLimit(index int) []IPriority {
	newHeap := Heap{}
	newHeap.heap = h.CopyAsArray()
	newHeap.size = h.size
	newHeap.capacity = h.capacity
	newHeap.flag = h.flag
	array := make([]IPriority, index)
	i := 0
	//fmt.Printf("newheap add is %x ,hHeap addr is %x  \n" ,newHeap,h)
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

func (h Heap) CopyAsMapTopPriceLimit(limit int) map[string]float64 {
	orderMap := make(map[string]float64)
	keyCount := 0
	//var heap Heap = Heap(h)
	var i = 0
	var item IPriority
	for {
		item, i = h.Next(i)
		if i == -1 {
			break
		}
		switch t := item.(type) {
		case order.Order:
			//fmt.Printf("type  %v \n", t.Price)
			num, ok := orderMap[fmt.Sprintf("%f", t.Price)]
			if ok {
				//fmt.Println(t.Price)
				orderMap[fmt.Sprintf("%f", t.Price)] = num + t.Num
			} else {

				orderMap[fmt.Sprintf("%f", t.Price)] = t.Num
				keyCount++

				if keyCount > limit {
					delete(orderMap, fmt.Sprintf("%f", t.Price))
					break
				}
			}
		default:
			fmt.Printf("unknown type %T \n", t)
		}

	}

	return orderMap
}

/*
func FormattedJson(obj *map[string]float64) string {
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
}*/

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

func (h Heap) SortByPrice(m map[string]float64) {
	var keys []float64
	var values []float64
	for key, value := range m {
		v, _ := strconv.ParseFloat(key, 64)
		keys = append(keys, v)
		values = append(values, value)
	}
	if h.isMaxRootHeap() {
		//sort.Sort(sort.StringSlice(keys))
		sort.Sort(sort.Reverse(sort.Float64Slice(keys)))
	} else {
		sort.Sort(sort.Reverse(sort.Float64Slice(keys)))
		//sort.Sort(sort.StringSlice(keys))
	}

	for i, key := range keys {
		fmt.Printf(" %v,  %v\n", key, values[i])
	}
}
