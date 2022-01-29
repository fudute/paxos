package store

import (
	"sync"
)

type Status int

const (
	Idle    Status = 0
	Running Status = 1
	Chosen  Status = 2
)

type Picker struct {
	mu           sync.Mutex
	window       []Status
	smallestIdle int
	off          int
}

func NewPicker() Picker {
	return Picker{off: 1}
}

func (p *Picker) Pick() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	ind := p.smallestIdle - p.off
	for ind < len(p.window) && p.window[ind] != Idle {
		ind++
	}

	if ind == len(p.window) {
		p.window = append(p.window, make([]Status, len(p.window)+1)...)
	}

	p.window[ind] = Running
	p.smallestIdle = ind + p.off

	return p.smallestIdle
}

func (p *Picker) Put(i int) {
	if i < p.off {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	for i-p.off >= len(p.window) {
		p.window = append(p.window, make([]Status, len(p.window)+1)...)
	}

	p.window[i-p.off] = Chosen

	for i := 0; i < len(p.window); i++ {
		if p.window[i] == Chosen {
			if p.off == p.smallestIdle {
				p.smallestIdle++
			}
			p.off++
			p.window = p.window[1:]
			continue
		}
		break
	}
}
