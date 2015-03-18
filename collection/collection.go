package collection

import (
	"github.com/boomlinde/acidforth/machine"
	"github.com/boomlinde/acidforth/machine/stack"
	"sync"
)

type Collection struct {
	Mutex   sync.Mutex
	tickers []func()
	Machine *machine.Machine
	out1    float32
	out2    float32
}

func (c *Collection) Register(ticker func()) {
	c.tickers = append(c.tickers, ticker)
}

func (c *Collection) Callback(buf [][]float32) {
	c.Mutex.Lock()
	for i := range buf[0] {
		for _, t := range c.tickers {
			t()
		}
		c.out1 = 0
		c.out2 = 0
		c.Machine.Run()
		buf[0][i] = c.out1
		buf[1][i] = c.out2
	}
	c.Mutex.Unlock()
}

func NewCollection() *Collection {
	col := &Collection{
		tickers: make([]func(), 0),
		Machine: machine.NewMachine(),
	}
	col.Machine.Register(">out1", func(s *stack.Stack) { col.out1 = float32(s.Pop()) })
	col.Machine.Register(">out2", func(s *stack.Stack) { col.out2 = float32(s.Pop()) })
	col.Machine.Compile(machine.TokenizeString(""))
	return col
}
