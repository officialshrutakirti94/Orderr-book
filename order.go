package main

import (
	"container/heap"
	"fmt"
)

// Order represents a buy or sell order
type Order struct {
	ID       int
	Price    int
	Quantity int
	IsBuy    bool
}

// OrderBook represents the order book
type OrderBook struct {
	buyHeap  *orderHeap
	sellHeap *orderHeap
}

// NewOrderBook creates a new order book
func NewOrderBook() *OrderBook {
	return &OrderBook{
		buyHeap:  &orderHeap{buy: true},
		sellHeap: &orderHeap{},
	}
}

// AddOrder adds a new order to the order book
func (ob *OrderBook) AddOrder(order Order) {
	if order.IsBuy {
		heap.Push(ob.buyHeap, order)
	} else {
		heap.Push(ob.sellHeap, order)
	}
}

// MatchOrders matches buy and sell orders
func (ob *OrderBook) MatchOrders() {
	for ob.buyHeap.Len() > 0 && ob.sellHeap.Len() > 0 {
		buyOrder := ob.buyHeap.Peek()
		sellOrder := ob.sellHeap.Peek()

		if buyOrder.Price < sellOrder.Price {
			break // no more matching possible
		}

		quantity := min(buyOrder.Quantity, sellOrder.Quantity)
		fmt.Printf("Matched order: Buy %d @ %d with Sell %d @ %d\n", quantity, buyOrder.Price, quantity, sellOrder.Price)

		buyOrder.Quantity -= quantity
		sellOrder.Quantity -= quantity

		if buyOrder.Quantity == 0 {
			heap.Pop(ob.buyHeap)
		}
		if sellOrder.Quantity == 0 {
			heap.Pop(ob.sellHeap)
		}
	}
}

// orderHeap is a min-heap implementation for orders
type orderHeap struct {
	orders []Order
	buy    bool
}

func (oh orderHeap) Len() int { return len(oh.orders) }

func (oh orderHeap) Less(i, j int) bool {
	if oh.buy {
		return oh.orders[i].Price >= oh.orders[j].Price
	}
	return oh.orders[i].Price <= oh.orders[j].Price
}

func (oh orderHeap) Swap(i, j int) { oh.orders[i], oh.orders[j] = oh.orders[j], oh.orders[i] }

func (oh *orderHeap) Push(x interface{}) { oh.orders = append(oh.orders, x.(Order)) }

func (oh *orderHeap) Pop() interface{} {
	old := oh.orders
	n := len(old)
	x := old[n-1]
	oh.orders = old[0 : n-1]
	return x
}

func (oh *orderHeap) Peek() Order {
	return oh.orders[0]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	orderBook := NewOrderBook()

	orderBook.AddOrder(Order{ID: 1, Price: 100, Quantity: 5, IsBuy: true})
	orderBook.AddOrder(Order{ID: 2, Price: 110, Quantity: 7, IsBuy: true})
	orderBook.AddOrder(Order{ID: 3, Price: 105, Quantity: 3, IsBuy: false})
	orderBook.AddOrder(Order{ID: 4, Price: 98, Quantity: 8, IsBuy: true})
	orderBook.AddOrder(Order{ID: 5, Price: 106, Quantity: 4, IsBuy: false})

	orderBook.MatchOrders()
}
