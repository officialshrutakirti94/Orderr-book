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
	buyHeap     *orderHeap
	sellHeap    *orderHeap
	MarketPrice int // Indicative market price
}

// NewOrderBook creates a new order book
func NewOrderBook() *OrderBook {
	return &OrderBook{
		buyHeap:     &orderHeap{buy: true},
		sellHeap:    &orderHeap{},
		MarketPrice: 0, // Initialize market price to zero
	}
}

// AddOrder adds a new order to the order book
// AddOrder adds a new order to the order book
func (ob *OrderBook) AddOrder(order Order) {
	if order.IsBuy {
		heap.Push(ob.buyHeap, order)
	} else {
		heap.Push(ob.sellHeap, order)
	}
	fmt.Println("Added order:", order)
	ob.UpdateMarketPrice() // Update market price after adding order
}

// MatchOrders matches buy and sell orders
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
			ob.UpdateMarketPrice() // Update market price after removing a sell order
		}
	}

	// Update market price after matching orders
	ob.UpdateMarketPrice()
}

// UpdateMarketPrice updates the market price based on the current bid-ask spread
// UpdateMarketPrice updates the market price based on the current bid-ask spread
func (ob *OrderBook) UpdateMarketPrice() {
	// Initialize best bid and best ask prices
	bestBid := 0
	bestAsk := 0

	// Check if there are buy orders (bids) in the order book
	if ob.buyHeap.Len() > 0 {
		bestBid = ob.buyHeap.Peek().Price
	}

	// Check if there are sell orders (asks) in the order book
	if ob.sellHeap.Len() > 0 {
		bestAsk = ob.sellHeap.Peek().Price
	}

	// Debug print statements
	fmt.Println("Best Bid:", bestBid)
	fmt.Println("Best Ask:", bestAsk)

	// Calculate the mid-price (indicative market price)
	// as the average of the highest bid and lowest ask prices
	if bestBid != 0 && bestAsk != 0 {
		ob.MarketPrice = (bestBid + bestAsk) / 2
	} else if bestBid != 0 {
		ob.MarketPrice = bestBid
	} else if bestAsk != 0 {
		ob.MarketPrice = bestAsk
	} else {
		ob.MarketPrice = 0 // If there are no bids or asks, set market price to 0
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
		return oh.orders[i].Price > oh.orders[j].Price // Higher price is considered "less" for buy orders
	}
	return oh.orders[i].Price < oh.orders[j].Price // Lower price is considered "less" for sell orders
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
	if oh.Len() == 0 {
		return Order{} // Return an empty order if the heap is empty
	}
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

	// Add some sample orders
	orderBook.AddOrder(Order{ID: 1, Price: 100, Quantity: 5, IsBuy: true})
	orderBook.AddOrder(Order{ID: 2, Price: 110, Quantity: 7, IsBuy: true})
	orderBook.AddOrder(Order{ID: 3, Price: 105, Quantity: 3, IsBuy: false})
	orderBook.AddOrder(Order{ID: 4, Price: 98, Quantity: 8, IsBuy: true})
	orderBook.AddOrder(Order{ID: 5, Price: 106, Quantity: 4, IsBuy: false})

	// Match orders and update market price
	orderBook.MatchOrders()

	// Print the indicative market price
	fmt.Printf("Indicative Market Price: %d\n", orderBook.MarketPrice)
}
