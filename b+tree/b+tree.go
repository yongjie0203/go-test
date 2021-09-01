package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type BPlusTree struct {
	//
	MaxDegree int

	size int

	root *IndexNode
}

type DataNode struct {
	data []float64
}

type IndexNode struct {
	childNode []*IndexNode

	DataNode
}

func (t *BPlusTree) maxDataSize() int {
	return t.MaxDegree
}

func InitTree() *BPlusTree {
	maxDegree := 2
	size := 0
	root := &IndexNode{DataNode: DataNode{data: make([]float64, maxDegree)}}
	return &BPlusTree{MaxDegree: maxDegree, size: size, root: root}
}

func (t *BPlusTree) insert(item float64) {

	if t.size == 0 {
		data := make([]float64, t.maxDataSize())
		data = append(data, item)
		indexNode := IndexNode{DataNode: DataNode{data: data}}
		t.root = &indexNode
		t.size++
	}
}

func (t *BPlusTree) insertRoot(item float64) {

	if t.size == 0 {
		data := make([]float64, t.maxDataSize())
		data = append(data, item)
		indexNode := IndexNode{DataNode: DataNode{data: data}}
		t.root = &indexNode
		t.size++
	}

	if t.size < t.MaxDegree {

	}
}

func (t *BPlusTree) sort(data []float64) {
	for i := 0; i < t.MaxDegree; i++ {

	}
}

func (t *BPlusTree) swap(data []interface{}, i, j int) {
	data[i], data[j] = data[j], data[i]
}

func main() {
	tree := InitTree()

	fmt.Printf("%d,%s \n", 1, FormattedJson(tree))

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
