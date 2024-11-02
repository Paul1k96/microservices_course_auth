package closer

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = New()

// Add adds `func() error` callback to the globalCloser
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait ...
func Wait() {
	globalCloser.Wait()
}

// CloseAll ...
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer ...
type Closer struct {
	mu         sync.Mutex
	once       sync.Once
	logger     *slog.Logger
	done       chan struct{}
	funcsServe []func() error
}

// New returns new Closer, if []os.Signal is specified Closer will automatically
// call CloseAll when one of signals is received from OS
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

// Add func to closer
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcsServe = append(c.funcsServe, f...)
	c.mu.Unlock()
}

// Wait blocks until all closer functions are done
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll calls all added functions and blocks until all are done
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcsServe
		c.funcsServe = nil
		c.mu.Unlock()

		// call all Closer funcs async
		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				c.logger.Error("error returned from Closer")
			}
		}
	})
}
