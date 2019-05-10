package models

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

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
		panic(ErrSide)
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
		panic(ErrSide)
	}
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
		return ErrSide
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

// OrderCache is a order input & output channel
type OrderCache struct {
	inputs            chan *Order
	results           chan *Order
	orderCache        map[string]*Order
	inflightCache     map[string]*Order
	maxInflightOrders int
	orderRate         float64
}

func (cache *OrderCache) requireToken(timeout time.Duration) <-chan error {
	return make(<-chan error)
}

// Put order into order cache, it's go routing safe
func (cache *OrderCache) Put(ord *Order, timeout time.Duration) error {
	if int64(timeout) > 0 {
		select {
		case cache.inputs <- ord:
			return nil
		case <-time.After(timeout):
			return fmt.Errorf("put order timeout: %v", timeout)
		}
	} else {
		cache.inputs <- ord
		return nil
	}
}

// PutResult puts order result into cache
func (cache *OrderCache) PutResult(ord *ngerest.Order) {
	// todo: inflight order handle
	converted := ConvertOrder(ord)

	if converted != nil {
		cache.results <- converted
	}
}

// GetResults to get order results channl
func (cache *OrderCache) GetResults() <-chan *Order { return cache.results }

// CloseResults to close order results channel
func (cache *OrderCache) CloseResults() { close(cache.results) }

// NewOrderCache to make new order cache
func NewOrderCache() *OrderCache {
	cache := OrderCache{
		inputs:        make(chan *Order),
		results:       make(chan *Order),
		orderCache:    make(map[string]*Order),
		inflightCache: make(map[string]*Order),
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
