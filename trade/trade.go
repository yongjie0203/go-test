package trade

import (
	"fmt"
	"log"

	"github.com/yongjie0203/go-trade-core/heap"
	"github.com/yongjie0203/go-trade-core/order"
)

func peekOrder(heap *heap.Heap) order.Order {
	var nilOrder order.Order
	if heap.Size() > 0 {
		top, _ := heap.Peek()
		switch order := top.(type) {
		case order.Order:
			//fmt.Printf("peek addr : %x  order addr : %x \n", &top, order)
			return order
		default:
			fmt.Printf("unknown peekOrder type %T \n", order)
			return nilOrder
		}
	} else {
		return nilOrder
	}

}

func pollOrder(heap *heap.Heap) order.Order {
	var nilOrder order.Order
	if heap.Size() > 0 {
		top, _ := heap.Poll()
		switch order := top.(type) {
		case order.Order:
			return order
		default:
			fmt.Printf("unknown pollOrder type %T \n", order)
			return nilOrder
		}
	} else {
		return nilOrder
	}

}

var Count int = 0

func onTransaction(sell, buy order.Order, transactionNum float64) {

	Count++
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

func Transaction(sell, buy *heap.Heap) {

	for {

		/*sell.add.Lock()
		buy.add.Lock()
		defer sell.add.Unlock()
		defer buy.add.Unlock()*/
		if sell.Size() > 0 && buy.Size() > 0 {

			deal(sell, buy)

		} else {
			//time.Sleep(time.Nanosecond * 10000)
			//runtime.Gosched()
			//select {}
			//fmt.Print("all order is deal closed \n")
		}
		//return false

	}
}

func deal(sell, buy *heap.Heap) {

	sell.Lock()
	defer sell.Unlock()
	buy.Lock()
	defer buy.Unlock()

	sellTopOrder := peekOrder(sell)
	buyTopOrder := peekOrder(buy)
	//fmt.Printf(" sell.Price : %f buy.Price : %f \n sell.Num:%f ,buy.Num:%f \n\n\n sell.heap:%v\n buy.heap:%v\n", sellTopOrder.Price, buyTopOrder.Price, sellTopOrder.Num, buyTopOrder.Num, heap.FormattedJson(sell.CopyAsSortArrayLimit(10)), heap.FormattedJson(buy.CopyAsSortArrayLimit(10)))

	if sellTopOrder.Price <= buyTopOrder.Price {
		//transactionPrice = int(sellTopOrder.Price)
		var transactionNum float64
		if sellTopOrder.Num <= buyTopOrder.Num {
			transactionNum = sellTopOrder.Num
		} else {
			transactionNum = buyTopOrder.Num
		}

		onTransaction(sellTopOrder, buyTopOrder, transactionNum)

		sellTopOrder.Num = sellTopOrder.Num - transactionNum
		buyTopOrder.Num = buyTopOrder.Num - transactionNum

		func() {

			scallback := func() {
				fmt.Printf("Head is miss ,reInsert %f \n", sellTopOrder.GetPriority())
				//go sell.Insert(sellTopOrder)
			}
			sell.Head(sellTopOrder, scallback)
			if sellTopOrder.Num <= 0 {

				//sell.poll()
				pollSell := pollOrder(sell)
				//test code

				if pollSell.OID != sellTopOrder.OID {
					//log.Printf("sell Oid not match %v %v \n", pollSell.OID, sellTopOrder.OID)
					//break
					//return true
					//sell.Insert(pollSell)
					panic(fmt.Sprintf("sell Oid not match %v %v \n", pollSell.OID, sellTopOrder.OID))
				}
				if pollSell.Num != sellTopOrder.Num {
					//log.Printf("sell num not match %v %v \n", pollSell.Num, sellTopOrder.Num)
					//return true
					panic(fmt.Sprintf("sell num not match %v %v \n", pollSell.Num, sellTopOrder.Num))
				}
			}

		}()

		func() {

			bcallback := func() {
				fmt.Printf("Head is miss ,reInsert %f \n", buyTopOrder.GetPriority())
				//go buy.Insert(buyTopOrder)
			}
			buy.Head(buyTopOrder, bcallback)

			if buyTopOrder.Num <= 0 {

				//buy.poll()
				pollBuy := pollOrder(buy)
				//test code

				if pollBuy.OID != buyTopOrder.OID {
					//log.Printf("buy Oid not match %v %v \n", pollBuy.OID, buyTopOrder.OID)
					//sell.Insert(pollBuy)
					panic(fmt.Sprintf("buy Oid not match %v %v \n", pollBuy.OID, buyTopOrder.OID))
					//break
					//return true
				}
				if pollBuy.Num != buyTopOrder.Num {
					//log.Printf("buy Num not match %v %v \n", pollBuy.Num, buyTopOrder.Num)
					panic(fmt.Sprintf("buy Num not match %v %v \n", pollBuy.Num, buyTopOrder.Num))
					//break
					//return true
				}

			}

		}()

		//m := buy.CopyAsMapTopPriceLimit(5)
		//fmt.Printf("\n buy list %s \n", heap.FormattedJson(&m))
		//m = sell.CopyAsMapTopPriceLimit(5)
		//fmt.Printf("\n sell list %s \n", heap.FormattedJson(&m))
		//fmt.Printf(" sell.num : %f buy.num : %f \n", sellTopOrder.Num, buyTopOrder.Num)
		//fmt.Printf("buy heap sise is %d sell heap size is %d\n", buy.Size(), sell.Size())

	} else {
		//time.Sleep(time.Nanosecond * 10000)
	}
}
