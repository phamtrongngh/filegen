package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var count int
var prefix string
var dir string
var size int64
var help bool

func init() {
	flag.IntVar(&count, "n", 1, "The number of files will be generated")
	flag.StringVar(&prefix, "p", "file", "The prefix of files will be generated")
	flag.Int64Var(&size, "s", 1024, "The size of files in bytes")
	flag.StringVar(&dir, "d", ".", "The directory where files will be generated")
	flag.Parse()
}

func createFile(filename string, size int64) error {
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("File %s will not be generated because it exists\n", filename)
		return nil
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	rf, err := os.Open("/dev/random")
	if err != nil {
		return err
	}
	defer rf.Close()

	_, err = io.CopyN(f, rf, size)
	if err != nil {
		return err
	}
	fmt.Printf("File %s has been generated\n", filename)
	return nil
}

func main() {
	if help {
		flag.Usage()
		return
	}

	d, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	f, err := d.Stat()
	if err != nil {
		panic(err)
	}

	if !f.IsDir() {
		fmt.Println("path is not a directory")
		return
	}

	wg := sync.WaitGroup{}
	for i := 0; i < count; i++ {
		wg.Add(1)
		filename := filepath.Join(dir, fmt.Sprintf("%s-%d", prefix, i))
		go func(filename string) {
			defer wg.Done()
			if err := createFile(filename, size); err != nil {
				fmt.Println(err)
				return
			}
		}(filename)
	}

	wg.Wait()
}
