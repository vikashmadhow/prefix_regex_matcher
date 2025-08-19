package regex

import (
	"testing"
)

func TestGenerator_Next1(t *testing.T) {
	r := NewRegex("a*")
	for i := 0; i < 10; i++ {
		println(r.Generate())
	}
}

func TestGenerator_Next2(t *testing.T) {
	r := NewRegex("a+")
	for i := 0; i < 10; i++ {
		println(r.Generate())
	}
}

func TestGenerator_Next3(t *testing.T) {
	r := NewRegex("([A-Z][0-9]){3}|([A-Z][0-9][A-Z] [0-9][A-Z][0-9])")
	for i := 0; i < 10; i++ {
		println(r.Generate())
	}
}

func TestGenerator_Next4(t *testing.T) {
	r := NewRegex("\\d{4}-[A-Z]+")
	for i := 0; i < 10; i++ {
		println(r.Generate())
	}
}

func TestGenerator_Next5(t *testing.T) {
	r := NewRegex("[A-Z]{5}-[^A-Za-Z!#$#$#$#]{10}")
	for i := 0; i < 10; i++ {
		println(r.Generate())
	}
}
