package iter

import (
	"errors"
	"sync"
)

const (
	MinChunkSize uint = 256
	MaxChunkSize uint = 4096
)

type ParIter struct {
	errors []error
	mu     sync.Mutex

	wg    *sync.WaitGroup
	errCh chan error
	done  chan struct{}
}

func NewParIter() *ParIter {
	p := &ParIter{
		errors: []error{},
		wg:     &sync.WaitGroup{},
		errCh:  make(chan error),
		done:   make(chan struct{}),
	}

	go func() {
		for err := range p.errCh {
			p.mu.Lock()
			p.errors = append(p.errors, err)
			p.mu.Unlock()
		}
		close(p.done)
	}()

	return p
}

func (p *ParIter) Wait() error {
	p.wg.Wait()

	close(p.errCh)
	<-p.done

	p.mu.Lock()
	defer p.mu.Unlock()

	return errors.Join(p.errors...)
}

type ParIterOptions struct {
	MaxChunkSize uint
	MinChunkSize uint
}

func DefaultParIterOptions() *ParIterOptions {
	return &ParIterOptions{
		MaxChunkSize: MaxChunkSize,
		MinChunkSize: MinChunkSize,
	}
}

func WithMinChunkSize(size uint) ParIterOptions {
	return ParIterOptions{
		MinChunkSize: size,
	}
}

func WithMaxChunkSize(size uint) ParIterOptions {
	return ParIterOptions{
		MaxChunkSize: size,
	}
}

func BuildParIterOptions(opts []ParIterOptions) ParIterOptions {
	o := DefaultParIterOptions()

	for _, opt := range opts {
		if opt.MaxChunkSize > 0 {
			o.MaxChunkSize = opt.MaxChunkSize
		}
		if opt.MinChunkSize > 0 {
			o.MinChunkSize = opt.MinChunkSize
		}
	}

	return *o
}
