package share

import "runtime"

func WorkersCount() int {
	return runtime.NumCPU()
}

func ChunkSize(len int) int {
	workersCount := WorkersCount()

	targetChunks := workersCount * 4
	size := len / targetChunks

	if size < 1 {
		return 1
	}

	return size
}

func ChunkCount(len, chunkSize int) int {
	return (len + chunkSize - 1) / chunkSize
}

func ChunkIndexes(chunkIdx, len, chunkSize int) (int, int) {
	start := chunkIdx * chunkSize
	end := start + chunkSize

	if end > len {
		end = len
	}
	return start, end
}
