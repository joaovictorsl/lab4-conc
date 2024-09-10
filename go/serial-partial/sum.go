package main

import (
	"fmt"
	"os"
)

type fileSum struct {
	path     string
	totalSum int64
	chunks   map[int64]int
}

func sumChunks(filePath string, chunkSize int) (*fileSum, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fSum := &fileSum{filePath, 0, make(map[int64]int, 0)}

	var bytesRead int
	for bytesRead = 0; bytesRead < len(raw); {
		start := bytesRead
		end := start + chunkSize
		if end > len(raw) {
			end = len(raw)
		}

		partialSum := int64(0)
		for _, b := range raw[bytesRead:end] {
			partialSum += int64(b)
		}

		fSum.totalSum += partialSum
		fSum.chunks[partialSum]++
		bytesRead += end - start
	}

	return fSum, nil
}

// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	chunkSize := 2
	var totalSum int64
	sums := make(map[int64][]*fileSum)
	for _, path := range os.Args[1:] {
		fSum, err := sumChunks(path, chunkSize)

		if err != nil {
			continue
		}

		totalSum += fSum.totalSum
		fmt.Println(fSum.chunks)

		sums[fSum.totalSum] = append(sums[fSum.totalSum], fSum)
	}

	fmt.Println("totalSum", totalSum)

	for sum, files := range sums {
		if len(files) > 1 {
			fileName := make([]string, 0)
			for _, f := range files {
				fileName = append(fileName, f.path)
			}

			fmt.Printf("Sum %d: %v\n", sum, fileName)
		}
	}

}

/*
func checkSimilarity(f1 *fileSum, f2 *fileSum) float64 {
	matches := 0
	total := 0

	for v, q := range f1.chunks {
		q2, ok := f2.chunks[v]
		if !ok {
			q2 = 0
		}

		qMatch := 0
		if q < q2 {
			qMatch = q
		} else {
			qMatch = q2
		}

		matches += qMatch
	}

	return float64(matches) / total
}
*/
