package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/frozenpine/ngecli/models"
)

// ReadLine read line from io.Reader interface
func ReadLine(prompt string, src io.Reader) string {
	if prompt == "" {
		prompt = "Please input: "
	}

	if src == nil {
		src = os.Stdin
	}

	reader := bufio.NewReader(src)

	fmt.Print(prompt)

	text, _ := reader.ReadString('\n')

	return strings.TrimRight(text, "\r\n")
}

// CheckSymbol validate order symbol
func CheckSymbol(symbol string) error {
	if symbol == "" {
		return models.ErrSymbol
	}

	return nil
}

// CheckPrice validate order price
func CheckPrice(price float64) error {
	if price <= 0 {
		return models.ErrPrice
	}

	return nil
}

// CheckQuantity validate order quantity
func CheckQuantity(qty int64) error {
	if qty == 0 {
		return models.ErrQuantity
	}

	return nil
}

// MatchSide match side with quantity
func MatchSide(side *models.OrderSide, qty int64) error {
	switch *side {
	case models.Buy:
		if qty < 0 {
			return models.ErrMissMatchQtySide
		}
	case "":
		if qty > 0 {
			*side = models.Buy
		} else {
			*side = models.Sell
		}
	default:
		return models.ErrSide
	}

	return nil
}
