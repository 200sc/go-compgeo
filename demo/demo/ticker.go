package demo

import "time"

// A DynamicTicker is a ticker which can
// be sent signals in the form of durations to
// change how often it ticks.
type DynamicTicker struct {
	ticker    *time.Ticker
	ch        chan time.Time
	resetCh   chan *time.Ticker
	forceTick chan bool
}

// NewDynamicTicker returns a null-initialized
// dynamic ticker
func NewDynamicTicker() *DynamicTicker {
	ch := make(chan time.Time)
	resetCh := make(chan *time.Ticker)
	forceTick := make(chan bool)
	dt := &DynamicTicker{
		// Please do not leave the application running
		// for a thousand hours without clicking on
		// the visualization knub, or else your next
		// visualization animation might skip a frame!
		// (We need -some- ticker defined or else
		// the program will crash in the following
		// routine on a nil pointer)
		ticker:    time.NewTicker(1000 * time.Hour),
		ch:        ch,
		resetCh:   resetCh,
		forceTick: forceTick,
	}
	go func(dt *DynamicTicker) {
		for {
			select {
			case v := <-dt.ticker.C:
				select {
				case <-dt.forceTick:
					continue
				case dt.ch <- v:
				case ticker := <-dt.resetCh:
					dt.ticker.Stop()
					dt.ticker = ticker
				}
			case ticker := <-dt.resetCh:
				dt.ticker.Stop()
				dt.ticker = ticker
			case <-dt.forceTick:
				select {
				case <-dt.forceTick:
					continue
				case dt.ch <- time.Time{}:
				}
			}
		}
	}(dt)
	return dt
}

// SetTick changes the rate at which a dynamic ticker
// ticks
func (dt *DynamicTicker) SetTick(d time.Duration) {
	dt.resetCh <- time.NewTicker(d)
}

// Step will force the dynamic ticker to tick, once.
// If the forced tick is not received, multiple calls
// to step will do nothing.
func (dt *DynamicTicker) Step() {
	select {
	case dt.forceTick <- true:
	default:
	}
}
