// Package iter provides utilities for parallel iteration and
// chunk-based concurrent processing.
package iter

import (
	"errors"
	"sync"
)

const (
	// MinChunkSize is the minimum chunk size for
	// parallel processing.
	MinChunkSize uint = 256

	// MinChunkSize is the maximum chunk size for
	// parallel processing.
	MaxChunkSize uint = 4096
)

// ParIter coordinates parallel job execution and error aggregation.
type ParIter struct {
	jobsWg *sync.WaitGroup
	errors []error

	errCh          chan error
	errDoneCh      chan struct{}
	jobsDoneCh     chan struct{}
	postJobsDoneCh chan struct{}

	doneOnce sync.Once
}

// NewParIter creates and initializes a new ParIter.
func NewParIter() *ParIter {
	p := &ParIter{
		errors:         []error{},
		jobsWg:         &sync.WaitGroup{},
		errDoneCh:      make(chan struct{}),
		jobsDoneCh:     make(chan struct{}),
		postJobsDoneCh: make(chan struct{}),
		errCh:          make(chan error),
		doneOnce:       sync.Once{},
	}

	go func() {
		for err := range p.errCh {
			p.errors = append(p.errors, err)
		}
		close(p.errDoneCh)
	}()

	return p
}

// Wait blocks until all jobs have completed and any
// post-job processing has finished.
// It returns a single error composed using errors.Join. If no
// errors were reported, it returns nil.
func (p *ParIter) Wait() error {
	p.doneOnce.Do(func() {
		<-p.jobsDoneCh
		<-p.postJobsDoneCh
		close(p.errCh)
		<-p.errDoneCh
	})

	return errors.Join(p.errors...)
}

// postJobsDone signals that all post-job processing has completed.
func (p *ParIter) postJobsDone() {
	close(p.postJobsDoneCh)
}

// startJobsDoneWatcher starts a waiter that closes jobsDoneCh
// when all queued jobs have completed.
func (p *ParIter) startJobsDoneWatcher() {
	go func() {
		p.jobsWg.Wait()
		close(p.jobsDoneCh)
	}()
}

// ParIterOptions configures behavior for parallel iteration.
type ParIterOptions struct {
	MaxChunkSize uint
	MinChunkSize uint
}

// Returns a ParIterOptions instance populated with
// default values.
func DefaultParIterOptions() *ParIterOptions {
	return &ParIterOptions{
		MaxChunkSize: MaxChunkSize,
		MinChunkSize: MinChunkSize,
	}
}

// Returns an option that overrides the minimum
// chunk size for parallel processing.
func WithMinChunkSize(size uint) ParIterOptions {
	return ParIterOptions{
		MinChunkSize: size,
	}
}

// Returns an option that overrides the maximum
// chunk size for parallel processing.
func WithMaxChunkSize(size uint) ParIterOptions {
	return ParIterOptions{
		MaxChunkSize: size,
	}
}

// Merges multiple ParIterOptions into a single
// configuration instance.
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
