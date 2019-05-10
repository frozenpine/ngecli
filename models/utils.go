package models

import (
	"fmt"

	"github.com/frozenpine/ngerest"
)

// LogError parse error & log it
func LogError(err error) {
	if swErr, ok := err.(ngerest.GenericSwaggerError); ok {
		fmt.Printf("%s: %s", swErr.Error(), string(swErr.Body()))
	} else {
		fmt.Println(err)
	}
}
