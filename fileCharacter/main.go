package main

import (
	"github.com/saintfish/chardet"
	//"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func handleFile(path string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	//if directory or isn't csv file
	if f.IsDir() || !strings.Contains(f.Name(), ".csv") {
		return nil
	}
	file, err := os.Open(f.Name())
	if err != nil {
		return err
	}

	defer file.Close()
	size, _ := io.ReadFull(file, buffer)
	input := buffer[:size]
	result, err := detector.DetectBest(input)
	if err != nil {
		return err
	}
	if result.Charset != "UTF-8" {
		fmt.Printf("file charset error:%s, %s\n", f.Name(), result.Charset)
	}
	return nil
}

func getFilelist(path string, f filepath.WalkFunc) {
	err := filepath.Walk(path, f)
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

var path1 string = "./"
var buffer = make([]byte, 32<<10)
var detector = chardet.NewTextDetector()

func main() {
	//flag.Parse()
	//path1 = flag.Arg(0)
	fmt.Println("current path1 args:", path1)
	getFilelist(path1, handleFile)
	time.Sleep(time.Hour)
}
