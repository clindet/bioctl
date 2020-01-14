package file

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"sync"
)

// CounterByNameSlice abstract func to conduct slice counter
func CounterByNameSlice(files *[]string,
	f func(filename string) (int64, error)) (*map[string]int64, int64, []error) {
	countMap := make(map[string]int64)
	errs := []error{}

	if len(*files) == 0 {
		return &countMap, 0, nil
	}
	var wg sync.WaitGroup
	var total int64
	var lock sync.Mutex
	wg.Add(len(*files))
	for i := range *files {
		go func(i int64) {
			c, err := f((*files)[i])
			if err != nil {
				errs = append(errs, err)
			} else {
				lock.Lock()
				countMap[(*files)[i]] = c
				lock.Unlock()
				total = total + c
			}
			wg.Done()
		}(int64(i))
	}
	wg.Wait()
	sslice := []string{}
	for key := range countMap {
		sslice = append(sslice, key)
	}
	sort.Strings(sslice)
	for _, v := range sslice {
		fmt.Printf("%d\t%s\n", countMap[v], v)
	}
	if len(sslice) > 0 {
		fmt.Printf("%d\ttotal\n", total)
	}
	return &countMap, total, errs
}

// LineCounter can be used to count file N lines.
// modifed from https://ask.csdn.net/questions/1009986?sort=comments_count
func LineCounter(r io.Reader) (int64, error) {
	var newLineChr byte
	var readSizeTmp int
	var readSize int64
	var err error
	var count int64
	buf := make([]byte, 1024)
	if runtime.GOOS == "linux" {
		newLineChr = '\n'
	} else if runtime.GOOS == "darwin" {
		newLineChr = '\r'
	} else {
		newLineChr = '\n'
	}

	for {
		readSizeTmp, err = r.Read(buf)
		readSize = int64(readSizeTmp)
		if err != nil {
			break
		}

		var buffPosition int64
		for {
			i := int64(bytes.IndexByte(buf[buffPosition:], newLineChr))
			if i == -1 || readSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
	}
	if readSize > 0 && count == 0 || count > 0 {
		count++
	}
	if err == io.EOF {
		return count - 1, nil
	}

	return count - 1, err
}

// LineCounterByName count lines from file
func LineCounterByName(filename string) (int64, error) {
	var err error
	var count int64
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return -1, err
	}
	var r *os.File
	if r, err = os.Open(filename); err != nil {
		return -1, err
	}
	count, err = LineCounter(r)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// LineCounterByNameSlice count lines from filename slice
// returns countMap, total , errs
func LineCounterByNameSlice(files []string) (*map[string]int64, int64, []error) {
	return CounterByNameSlice(&files, LineCounterByName)
}

// BytesCounter calculate the bytes of input stream
func BytesCounter(r io.Reader) (int64, error) {
	var readSizeTmp int
	var err error
	buf := make([]byte, 1024)
	var total int64
	for {
		readSizeTmp, err = r.Read(buf)
		total = total + int64(readSizeTmp)
		if err != nil {
			break
		}
		if err == io.EOF {
			return total, nil
		}
	}
	return total, err
}

// BytesCounterByName calculate the bytes of input stream
func BytesCounterByName(filename string) (int64, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return -1, err
	}
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return -1, err
	}
	return fileInfo.Size(), nil
}

// BytesCounterByNameSlice count bytes from filename slice
// returns countMap, total , errs
func BytesCounterByNameSlice(files []string) (*map[string]int64, int64, []error) {
	return CounterByNameSlice(&files, BytesCounterByName)
}
