package demo

import "time"

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
	}(dt)
	return dt
}

func (dt *dynamicTicker) SetTick(d time.Duration) {
	dt.resetCh <- time.NewTicker(d)
}
