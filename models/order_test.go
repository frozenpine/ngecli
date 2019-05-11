package models

import (
	"testing"
)

func TestOrderSide(t *testing.T) {
	side := new(OrderSide)

	if err := side.MatchSide(1); err != nil {
		t.Fatal(err)
	} else if *side != Buy {
		t.Fatal("auto config side failed.")
	}

	if err := side.MatchSide(-1); err == nil {
		t.Fatal("match side failed.")
	}

	sideValue := OrderSide("")

	if sideValue.MatchSide(-1); sideValue != Sell {
		t.Fatal("can not use value in MatchSide")
	} else {
		t.Log(sideValue)
	}
}
