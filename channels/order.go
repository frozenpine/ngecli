package channels

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/frozenpine/ngecli/models"
	"github.com/frozenpine/ngerest"
)

// OrderResultChan order result output channel
var OrderResultChan = make(chan *ngerest.Order)

var orderResultCache = make(map[string]*models.Order)

// ConvertOrder convert ngerest.Order to local model Order
func ConvertOrder(ord *ngerest.Order) *models.Order {
	if ord == nil {
		ErrChan <- errors.New("nil pointer in ConvertOrder function")
		return nil
	}

	var ordBuffer bytes.Buffer

	enc := gob.NewEncoder(&ordBuffer)
	dec := gob.NewDecoder(&ordBuffer)

	err := enc.Encode(ord)
	if err != nil {
		ErrChan <- err

		return nil
	}

	var converted models.Order

	err = dec.Decode(&converted)
	if err != nil {
		ErrChan <- err
		return nil
	}

	return &converted
}
