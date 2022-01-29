package store

import (
	"testing"
)

func TestPicker(t *testing.T) {
	p := &Picker{}
	p.Put(0)
	p.Put(1)
	p.Put(3)
	p.Put(4)
	if p.Pick() != 2 {
		t.Error()
	}
	if p.Pick() != 5 {
		t.Error()
	}
	p.Put(2)
	p.Put(6)
	if p.Pick() != 7 {
		t.Error()
	}
	for i := 0; i < 100; i++ {
		p.Put(p.Pick())
	}
}
