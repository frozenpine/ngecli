package common

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/frozenpine/ngerest"
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
		return ErrSymbol
	}

	return nil
}

// CheckPrice validate order price
func CheckPrice(price float64) error {
	if price <= 0 {
		return ErrPrice
	}

	return nil
}

// CheckQuantity validate order quantity
func CheckQuantity(qty int64) error {
	if qty == 0 {
		return ErrQuantity
	}

	return nil
}

// PrintError to auto parse err and print in console
func PrintError(prefix string, err error) {
	if swErr, ok := err.(ngerest.GenericSwaggerError); ok {
		fmt.Printf(
			prefix+": %s\n%s\n", swErr.Error(), string(swErr.Body()))
	} else {
		fmt.Printf(prefix+": %s\n", err.Error())
	}
}
