package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gocarina/gocsv"
)

func TestPassword(t *testing.T) {
	pass := NewPassword()

	passString := "yuanyang"

	fmt.Println(pass.Shadow(passString))

	if pass.Show() != passString {
		t.Fatal("shadow failed.")
	}

	shadowPass := "Kaia8b1g91tItzxpxgH1Syc4lZ2t02Gr9AHunwe4iSzO4EBVQkzLNuGMmwuW1Y0WRI302NiCVWZCIsjn+UPqgQKzNjWymCj3WIU4Ma8WH0gdunQJe8yStVBYX8RuF5SkroN8JArn8sBDSMLaHxJOgzB/+rBy8akVu61R0VoRS7Nsr5RTCCe/f3TwUlreRobcEo8hAhnWBjSvL+t4TXqpJrPJVr0nvYuvlhXM+iCzOKSWUdeYnmV29KiuBCl7GUa9TLgZgl/raHjAX45wRuKZI6tlffuJXA6SI7EKggO9Vh+UKUaJ3FnNSNRT2mv/CWBsC8jmm+MK1QwQzQiHEELuIw=="

	if err := pass.ShadowSet(shadowPass); err != nil {
		t.Fatal(err)
	}

	if pass.Show() != passString {
		t.Fatal("shadow failed.")
	}
}

func TestMarshal(t *testing.T) {
	csvContent := `
identity,password,api_key,api_secret
sonny.frozenpine@gmail.com,,ADa28R5s0dUfdn9W3STr,VQ9K28Rj9B35rmAe88VpCN04l1O3Hp1IpPe43y3U4MaSYzQKtijHj3om1dSCYeemagPX959pVj69Z5ESd9Q4T4rtT97h2j0k6Go
`

	var auths []*Authentication

	if err := gocsv.Unmarshal(strings.NewReader(csvContent), &auths); err != nil {
		t.Fatal(err)
	}

	for _, auth := range auths {
		if jsonBytes, err := json.Marshal(auth); err != nil {
			t.Fatal(err)
		} else {
			t.Log(string(jsonBytes))
		}
	}
}
