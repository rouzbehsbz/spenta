package types

import (
	"sync"

	"github.com/rouzbehsbz/spenta/internal/pool"
	"github.com/rouzbehsbz/spenta/internal/share"
)

type SliceParIter struct {
	wg *sync.WaitGroup
}

func NewSliceParIter[T comparable](slice *[]T, cb func(i int)) *SliceParIter {
	len, chunkSize, chunkCount := ChunkInfo(slice)

	wg := &sync.WaitGroup{}
	wg.Add(chunkCount)

	jobs := pool.NewSliceJobs(len, chunkCount, chunkSize, wg, func(i int) {
		cb(i)
	})

	go pool.SpentaPool().SendJobs(jobs)

	return &SliceParIter{
		wg: wg,
	}
}

func SliceParForEach[T comparable](slice *[]T, cb func(e T)) *SliceParIter {
	return NewSliceParIter[T](slice, func(i int) {
		cb((*slice)[i])
	})
}

func SliceParMap[T comparable](slice *[]T, cb func(e T) T) *SliceParIter {
	return NewSliceParIter[T](slice, func(i int) {
		(*slice)[i] = cb((*slice)[i])
	})
}

func SliceParFilter[T comparable](slice *[]T, cb func(e T) bool) *SliceParIter {
	return NewSliceParIter[T](slice, func(i int) {
		s := *slice
		j := 0
		for _, v := range s {
			if cb(v) {
				s[j] = v
				j++
			}
		}
		*slice = s[:j]
	})
}

func (p *SliceParIter) Done() {
	p.wg.Wait()
}

func ChunkInfo[T comparable](slice *[]T) (int, int, int) {
	len := len((*slice))
	chunkSize := share.ChunkSize(len)
	chunkCount := share.ChunkCount(len, chunkSize)

	return len, chunkSize, chunkCount
}
