package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string) (int, error) {
	data, err := readFile(filePath)
	if err != nil {
		return 0, err
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	return _sum, nil
}

// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	outCh := make(chan result, len(os.Args[1:]))

	for _, path := range os.Args[1:] {
		go sum_worker(path, outCh)
	}

	totalSum, sumMap := aggregator(len(os.Args[1:]), outCh)

	fmt.Println("totalSum:", totalSum)

	for sum, files := range sumMap {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}
}

type result struct {
	sum  int64
	path string
}

func sum_worker(path string, outCh chan result) {
	_sum, err := sum(path)

	if err != nil {
		return
	}

	outCh <- result{int64(_sum), path}
}

func aggregator(expected int, outCh chan result) (int64, map[int][]string) {
	var totalSum int64
	sumMap := make(map[int][]string)

	for i := 0; i < expected; i++ {
		res := <-outCh
		totalSum += res.sum
		sumMap[int(res.sum)] = append(sumMap[int(res.sum)], res.path)
	}
	close(outCh)

	return totalSum, sumMap
}
