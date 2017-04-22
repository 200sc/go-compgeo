package demo

import (
	"fmt"
	"time"
)

type dynamicTicker struct {
	ticker  *time.Ticker
	ch      chan time.Time
	resetCh chan *time.Ticker
}

func NewDynamicTicker() *dynamicTicker {
	ch := make(chan time.Time)
	resetCh := make(chan *time.Ticker)
	dt := &dynamicTicker{
		// Please do not leave the application running
		// for a thousand hours without clicking on
		// the visualization knub, or else your
		// visualization animation might skip a frame!
		// (We need -some- ticker defined or else
		// the program will crash in the following
		// routine on a nil pointer)
		ticker:  time.NewTicker(1000 * time.Hour),
		ch:      ch,
		resetCh: resetCh,
	}
	go func(dt *dynamicTicker) {
		for {
			fmt.Println("Waiting on ", dt.ticker.C, "or", dt.resetCh)
			select {
			case v := <-dt.ticker.C:
				select {
				case dt.ch <- v:
				case ticker := <-dt.resetCh:
					dt.ticker.Stop()
					dt.ticker = ticker
				}
			case ticker := <-dt.resetCh:
				dt.ticker.Stop()
				dt.ticker = ticker
			}
		}
		fmt.Println("Goroutine exited")
	}(dt)
	return dt
}

func (dt *dynamicTicker) SetTick(d time.Duration) {
	fmt.Println(dt.resetCh, "Sending on reset ch")
	dt.resetCh <- time.NewTicker(d)
}
