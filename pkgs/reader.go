package pkgs

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type Row struct {
	I int64
	S string
}

type FileReader struct {
	IgnoreLongRow bool
	err           error
	ch            chan Row
}

func (f *FileReader) Read(name string) (<-chan Row, error) {
	pf, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	var rd io.ReadCloser
	if filepath.Ext(name) == ".gz" {
		rd, err = gzip.NewReader(pf)
		if err != nil {
			return nil, err
		}
		fmt.Println("gz")
	} else {
		rd = pf
	}
	f.ch = make(chan Row, 1024)
	go f.read(rd)
	return f.ch, nil
}

func (f *FileReader) ReadConcurrentWithSkip(path string, isFirst bool, thread int, skip int64, fun func(i int64, s string) error) (int64, error) {
	ch, err := f.Read(path)
	if err != nil {
		return 0, err
	}
	var wg sync.WaitGroup
	var line int64
	for i := 0; i < thread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for s := range ch {
				if isFirst && s.I <= skip {
					continue
				}
				if e2 := fun(s.I, s.S); e2 != nil {
					f.err = e2
				} else {
					atomic.AddInt64(&line, 1)
				}
			}
		}()
	}
	wg.Wait()
	return line, f.err

}

func (f *FileReader) ReadConcurrent(path string, thread int, fun func(i int64, s string) error) (int64, error) {
	return f.ReadConcurrentWithSkip(path, false, thread, 0, fun)
}

func (f *FileReader) Error() error {
	return f.err
}

func (f *FileReader) read(pf io.ReadCloser) {
	defer pf.Close()
	rd := bufio.NewReader(pf)
	var (
		err error
		n   int64
	)
	for f.err == nil {
		n += 1
		line, err := rd.ReadBytes('\n')
		line = bytes.TrimRight(line, "\r\n")
		if err != nil {
			break
		}
		f.ch <- Row{I: n, S: string(line)}
	}
	if err != io.EOF {
		f.err = err
	}
	close(f.ch)
}
