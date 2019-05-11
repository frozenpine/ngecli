package models

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/frozenpine/ngecli/common"

	"github.com/frozenpine/ngerest"
)

// OrderSide order side
type OrderSide string

const (
	// Buy buy side
	Buy OrderSide = "Buy"
	// Sell sell side
	Sell OrderSide = "Sell"
)

// String order side string
func (s *OrderSide) String() string {
	return string(*s)
}

// Value order side value
func (s *OrderSide) Value() int64 {
	switch *s {
	case Buy:
		return 1
	case Sell:
		return -1
	default:
		panic(common.ErrSide)
	}
}

// Opposite get opposite order side
func (s *OrderSide) Opposite() OrderSide {
	switch *s {
	case Buy:
		return Sell
	case Sell:
		return Buy
	default:
		panic(common.ErrSide)
	}
}

// MatchSide match side with quantity
func (s *OrderSide) MatchSide(qty int64) error {
	switch *s {
	case Buy:
		if qty < 0 {
			return common.ErrMissMatchQtySide
		}
	case "":
		if qty > 0 {
			*s = Buy
		} else {
			*s = Sell
		}
	default:
		return common.ErrSide
	}

	return nil
}

// UnmarshalCSV unmarshal csv column to OrderSide
func (s *OrderSide) UnmarshalCSV(value string) error {
	return s.Set(value)
}

// MarshalCSV marshal to csv column
func (s *OrderSide) MarshalCSV() string {
	return (*s).String()
}

// UnmarshalJSON unmarshal from json string
func (s *OrderSide) UnmarshalJSON(data []byte) error {
	return s.Set(strings.Trim(string(data), "\""))
}

// MarshalJSON marshal to json string
func (s *OrderSide) MarshalJSON() ([]byte, error) {
	var buff bytes.Buffer
	buff.WriteString("\"" + (*s).String() + "\"")

	return buff.Bytes(), nil
}

// Set set OrderSide by string, if value is empty, default: Buy
func (s *OrderSide) Set(value string) error {
	switch value {
	case "Buy", "buy":
		*s = Buy
		return nil
	case "Sell", "sell":
		*s = Sell
		return nil
	case "":
		*s = Buy
		return nil
	default:
		return common.ErrSide
	}
}

// Type get order side type
func (s *OrderSide) Type() string {
	return s.String()
}

// Order order table
type Order struct {
	OrderID               string    `csv:"orderID" json:"orderID"`
	ClOrdID               string    `csv:"clOrdID,omitempty" json:"clOrdID,omitempty"`
	ClOrdLinkID           string    `csv:"clOrdLinkID,omitempty" json:"clOrdLinkID,omitempty"`
	Account               float32   `csv:"account,omitempty" json:"account,omitempty"`
	Symbol                string    `csv:"symbol,omitempty" json:"symbol,omitempty"`
	Side                  OrderSide `csv:"side,omitempty" json:"side,omitempty"`
	SimpleOrderQty        float64   `csv:"simpleOrderQty,omitempty" json:"simpleOrderQty,omitempty"`
	OrderQty              float32   `csv:"orderQty,omitempty" json:"orderQty,omitempty"`
	Price                 float64   `csv:"price,omitempty" json:"price,omitempty"`
	DisplayQty            float32   `csv:"displayQty,omitempty" json:"displayQty,omitempty"`
	StopPx                float64   `csv:"stopPx,omitempty" json:"stopPx,omitempty"`
	PegOffsetValue        float64   `csv:"pegOffsetValue,omitempty" json:"pegOffsetValue,omitempty"`
	PegPriceType          string    `csv:"pegPriceType,omitempty" json:"pegPriceType,omitempty"`
	Currency              string    `csv:"currency,omitempty" json:"currency,omitempty"`
	SettlCurrency         string    `csv:"settlCurrency,omitempty" json:"settlCurrency,omitempty"`
	OrdType               string    `csv:"ordType,omitempty" json:"ordType,omitempty"`
	TimeInForce           string    `csv:"timeInForce,omitempty" json:"timeInForce,omitempty"`
	ExecInst              string    `csv:"execInst,omitempty" json:"execInst,omitempty"`
	ContingencyType       string    `csv:"contingencyType,omitempty" json:"contingencyType,omitempty"`
	ExDestination         string    `csv:"exDestination,omitempty" json:"exDestination,omitempty"`
	OrdStatus             string    `csv:"ordStatus,omitempty" json:"ordStatus,omitempty"`
	Triggered             string    `csv:"triggered,omitempty" json:"triggered,omitempty"`
	WorkingIndicator      bool      `csv:"workingIndicator,omitempty" json:"workingIndicator,omitempty"`
	OrdRejReason          string    `csv:"ordRejReason,omitempty" json:"ordRejReason,omitempty"`
	SimpleLeavesQty       float64   `csv:"simpleLeavesQty,omitempty" json:"simpleLeavesQty,omitempty"`
	LeavesQty             float32   `csv:"leavesQty,omitempty" json:"leavesQty,omitempty"`
	SimpleCumQty          float64   `csv:"simpleCumQty,omitempty" json:"simpleCumQty,omitempty"`
	CumQty                float32   `csv:"cumQty,omitempty" json:"cumQty,omitempty"`
	AvgPx                 float64   `csv:"avgPx,omitempty" json:"avgPx,omitempty"`
	MultiLegReportingType string    `csv:"multiLegReportingType,omitempty" json:"multiLegReportingType,omitempty"`
	Text                  string    `csv:"text,omitempty" json:"text,omitempty"`
	TransactTime          JavaTime  `csv:"transactTime,omitempty" json:"transactTime,omitempty"`
	Timestamp             JavaTime  `csv:"timestamp,omitempty" json:"timestamp,omitempty"`
}

// IsClosed check if order is closed
func (ord *Order) IsClosed() bool {
	switch ord.OrdStatus {
	case "Filled":
		return true
	default:
		return false
	}
}

const (
	defaultInflightOrders int     = 5
	maxOrderRatePerUser   float64 = 1
	maxOrderRateTotal     float64 = 200
)

type clientCache struct {
	inQueue  map[string]*Order
	finished map[string]*Order
}

func (cc *clientCache) Queue(ord *Order) {
	cc.inQueue[ord.OrderID] = ord
}

func (cc *clientCache) Finish(ord *Order) {
	delete(cc.inQueue, ord.OrderID)

	cc.finished[ord.OrderID] = ord
}

func newClientCache() *clientCache {
	cache := clientCache{
		inQueue:  make(map[string]*Order),
		finished: make(map[string]*Order),
	}

	return &cache
}

// OrderCache is a order input & output channel
type OrderCache struct {
	inputs  chan *Order
	results chan *Order
	// orderCache OrderID as key
	orderCache     map[string]*Order
	inflightCache  map[string]*Order
	orderClientMap map[string]string
	// clientOrderCache client identity as key,
	clientOrderCache    map[string]*clientCache
	clientInflightQueue map[string]chan interface{}
	maxInflightOrders   int
	orderRate           float64
	tokenBucket         chan bool
}

func (cache *OrderCache) requireToken(timeChan <-chan time.Time) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			close(errChan)
		}()

		select {
		case <-cache.tokenBucket:
			return
		case <-timeChan:
			errChan <- common.ErrTokenInsufficient
		}
	}()

	return errChan
}

func (cache *OrderCache) checkInflight(
	id string, timeChan <-chan time.Time) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			close(errChan)
		}()

		queue, exist := cache.clientInflightQueue[id]
		if !exist {
			queue = make(chan interface{}, cache.maxInflightOrders)
			cache.clientInflightQueue[id] = queue
		}

		select {
		case queue <- nil:
			return
		case <-timeChan:
			errChan <- common.ErrInflightCheck
		}
	}()

	return errChan
}

func (cache *OrderCache) putOrder(id string, ord *Order) error {
	cache.inputs <- ord

	cache.inflightCache[ord.OrderID] = ord
	cache.orderCache[ord.OrderID] = ord
	cache.orderClientMap[ord.OrderID] = id

	clientCache, exist := cache.clientOrderCache[id]
	if !exist {
		clientCache = newClientCache()

		cache.clientOrderCache[id] = clientCache
	}

	clientCache.Queue(ord)

	return nil
}

func (cache *OrderCache) findClientIDByOrder(ord *Order) string {
	return cache.orderClientMap[ord.OrderID]
}

// PutOrder order into order cache, it's go routing safe
func (cache *OrderCache) PutOrder(
	id string, ord *Order, timeout time.Duration) error {
	var timeoutCh <-chan time.Time

	if int(timeout) > 0 {
		timeoutCh = time.After(timeout)
	} else {
		timeoutCh = make(<-chan time.Time)
	}

	if err := <-cache.checkInflight(id, timeoutCh); err != nil {
		return err
	}

	if err := <-cache.requireToken(timeoutCh); err != nil {
		return err
	}

	return cache.putOrder(id, ord)
}

// PutResult puts order result into cache
func (cache *OrderCache) PutResult(ord *ngerest.Order) {
	converted := ConvertOrder(ord)

	if converted == nil {
		jsonBytes, _ := json.Marshal(ord)
		fmt.Println("convert order failed, origin:", string(jsonBytes))
		return
	}

	cache.results <- converted

	clientID := cache.findClientIDByOrder(converted)

	if clientID == "" {
		fmt.Println("failed to find client id by order:", ord.OrderID)
		return
	}

	if converted.IsClosed() {
		if clientCache := cache.clientOrderCache[clientID]; clientCache != nil {
			clientCache.Finish(converted)
		} else {
			fmt.Println("client cache missingfor client:", clientID)
		}
	}

	if clientQueue := cache.clientInflightQueue[clientID]; clientQueue != nil {
		select {
		case <-clientQueue:
		default:
			fmt.Println("reduce inflight queue failed for client:", clientID)
		}
	} else {
		fmt.Println("inflight queue missing for client:", clientID)
	}
}

// GetResults to get order results channl
func (cache *OrderCache) GetResults() <-chan *Order { return cache.results }

// CloseResults to close order results channel
func (cache *OrderCache) CloseResults() { close(cache.results) }

// NewOrderCache to make new order cache
func NewOrderCache() *OrderCache {
	cache := OrderCache{
		inputs:              make(chan *Order),
		results:             make(chan *Order),
		orderCache:          make(map[string]*Order),
		inflightCache:       make(map[string]*Order),
		clientInflightQueue: make(map[string]chan interface{}),
		clientOrderCache:    make(map[string]*clientCache),
	}

	return &cache
}

// ConvertOrder convert ngerest.Order structure to local Order structure
func ConvertOrder(ori *ngerest.Order) *Order {
	var ordBuff bytes.Buffer

	enc := gob.NewEncoder(&ordBuff)
	dec := gob.NewDecoder(&ordBuff)

	err := enc.Encode(*ori)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var converted Order
	err = dec.Decode(&converted)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &converted
}
