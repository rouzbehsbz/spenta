package iter

import (
	"errors"
	"sync"
)

const MIN_CHUNK_SIZE uint = 512

type ParIter struct {
	errors []error

	wg    *sync.WaitGroup
	errCh chan error
}

func NewParIter(wg *sync.WaitGroup, errCh chan error) *ParIter {
	return &ParIter{
		errors: []error{},
		wg:     wg,
		errCh:  errCh,
	}
}

func (p *ParIter) Wait() error {
	p.wg.Wait()

	close(p.errCh)
	for err := range p.errCh {
		p.errors = append(p.errors, err)
	}

	return errors.Join(p.errors...)
}

type ParIterOptions struct {
	MinChunkSize uint
}

func DefaultParIterOptions() *ParIterOptions {
	return &ParIterOptions{
		MinChunkSize: MIN_CHUNK_SIZE,
	}
}

func WithMinChunkSize(size uint) ParIterOptions {
	return ParIterOptions{
		MinChunkSize: size,
	}
}

func BuildParIterOptions(opts []ParIterOptions) ParIterOptions {
	o := DefaultParIterOptions()

	for _, opt := range opts {
		o.MinChunkSize = opt.MinChunkSize
	}

	return *o
}

func ChunkSize(len int, minSize uint) int {
	if len > int(minSize) {
		len = (len + 1) / 2
		return ChunkSize(len, minSize)
	}

	return len
}

func ChunkCount(len, chunkSize int) int {
	return (len + chunkSize - 1) / chunkSize
}

func SliceChunk[V any](slice *[]V, minChunkSize uint) (int, int, int) {
	len := len((*slice))
	chunkSize := ChunkSize(len, minChunkSize)
	chunkCount := ChunkCount(len, chunkSize)

	return len, chunkSize, chunkCount
}
