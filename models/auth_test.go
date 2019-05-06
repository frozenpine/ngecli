package models

import (
	"fmt"
	"testing"
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
