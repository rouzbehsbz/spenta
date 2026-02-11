package internal

import "sync"

type ParSliceIter struct {
	Chunks int
	wg     *sync.WaitGroup
}

func SliceParIter[T comparable](slice *[]T, cb func(i int)) *ParSliceIter {
	len, chunkSize, chunkCount := ChunkInfo(slice)

	wg := &sync.WaitGroup{}
	wg.Add(chunkCount)

	jobs := NewSliceJobs(len, chunkCount, chunkSize, wg, func(i int) {
		cb(i)
	})

	go SpentaPool().SendJobs(jobs)

	return &ParSliceIter{
		Chunks: chunkSize,
		wg:     wg,
	}
}

func SliceParForEach[T comparable](slice *[]T, cb func(e T)) *ParSliceIter {
	return SliceParIter[T](slice, func(i int) {
		cb((*slice)[i])
	})
}

func SliceParMap[T comparable](slice *[]T, cb func(e T) T) *ParSliceIter {
	return SliceParIter[T](slice, func(i int) {
		(*slice)[i] = cb((*slice)[i])
	})
}

func (p *ParSliceIter) Done() {
	p.wg.Wait()
}

func ChunkInfo[T comparable](slice *[]T) (int, int, int) {
	len := len((*slice))
	chunkSize := ChunkSize(len)
	chunkCount := ChunkCount(len, chunkSize)

	return len, chunkSize, chunkCount
}
