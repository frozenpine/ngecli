package models

import (
	"bytes"
	"errors"
	"strings"
)

var (
	// ErrSide invalid side
	ErrSide = errors.New("side is either \"Buy\" or \"Sell\"")
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
func (s OrderSide) String() string {
	return string(s)
}

// Value order side value
func (s OrderSide) Value() int64 {
	switch s {
	case Buy:
		return 1
	case Sell:
		return -1
	default:
		panic(ErrSide)
	}
}

// Opposite get opposite order side
func (s OrderSide) Opposite() OrderSide {
	switch s {
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
	buff.WriteString((*s).String())

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

// Order order table
type Order struct {
	OrderID               string    `csv:"orderID" json:"orderID"`
	ClOrdID               string    `csv:"clOrdID,omitempty"`
	ClOrdLinkID           string    `csv:"clOrdLinkID,omitempty"`
	Account               float32   `csv:"account,omitempty"`
	Symbol                string    `csv:"symbol,omitempty"`
	Side                  OrderSide `csv:"side,omitempty"`
	SimpleOrderQty        float64   `csv:"simpleOrderQty,omitempty"`
	OrderQty              float32   `csv:"orderQty,omitempty"`
	Price                 float64   `csv:"price,omitempty"`
	DisplayQty            float32   `csv:"displayQty,omitempty"`
	StopPx                float64   `csv:"stopPx,omitempty"`
	PegOffsetValue        float64   `csv:"pegOffsetValue,omitempty"`
	PegPriceType          string    `csv:"pegPriceType,omitempty"`
	Currency              string    `csv:"currency,omitempty"`
	SettlCurrency         string    `csv:"settlCurrency,omitempty"`
	OrdType               string    `csv:"ordType,omitempty"`
	TimeInForce           string    `csv:"timeInForce,omitempty"`
	ExecInst              string    `csv:"execInst,omitempty"`
	ContingencyType       string    `csv:"contingencyType,omitempty"`
	ExDestination         string    `csv:"exDestination,omitempty"`
	OrdStatus             string    `csv:"ordStatus,omitempty"`
	Triggered             string    `csv:"triggered,omitempty"`
	WorkingIndicator      bool      `csv:"workingIndicator,omitempty"`
	OrdRejReason          string    `csv:"ordRejReason,omitempty"`
	SimpleLeavesQty       float64   `csv:"simpleLeavesQty,omitempty"`
	LeavesQty             float32   `csv:"leavesQty,omitempty"`
	SimpleCumQty          float64   `csv:"simpleCumQty,omitempty"`
	CumQty                float32   `csv:"cumQty,omitempty"`
	AvgPx                 float64   `csv:"avgPx,omitempty"`
	MultiLegReportingType string    `csv:"multiLegReportingType,omitempty"`
	Text                  string    `csv:"text,omitempty"`
	TransactTime          JavaTime  `csv:"transactTime,omitempty"`
	Timestamp             JavaTime  `csv:"timestamp,omitempty"`
}
