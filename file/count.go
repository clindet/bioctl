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

// LineCounter can be used to count file N lines.
// modifed from https://ask.csdn.net/questions/1009986?sort=comments_count
func LineCounter(r io.Reader) (int, error) {
	var newLineChr byte
	var readSize int
	var err error
	var count int

	buf := make([]byte, 1024)
	if runtime.GOOS == "linux" {
		newLineChr = '\n'
	} else if runtime.GOOS == "darwin" {
		newLineChr = '\r'
	} else {
		newLineChr = '\n'
	}

	for {
		readSize, err = r.Read(buf)
		if err != nil {
			break
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], newLineChr)
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

// LineCounterName count lines from file
func LineCounterName(filename string) (count int, err error) {
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

// CountLineNameSlice count lines from filename slice
// returns countMap, total , errs
func CountLineNameSlice(files []string) (*map[string]int, int, []error) {
	countMap := make(map[string]int)
	errs := []error{}

	if len(files) == 0 {
		return &countMap, 0, nil
	}
	var wg sync.WaitGroup
	var total int
	var lock sync.Mutex
	wg.Add(len(files))
	for i := range files {
		go func(i int) {
			c, err := LineCounterName(files[i])
			if err != nil {
				errs = append(errs, err)
			} else {
				lock.Lock()
				countMap[files[i]] = c
				lock.Unlock()
				total = total + c
			}
			wg.Done()
		}(i)
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
