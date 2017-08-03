package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type apkInfo struct {
	fileName   string
	modifyTime string
	size       int64
	diffSize   int64
}

var lowVersionMap = make(map[string]*apkInfo, 10000)
var highVersionMap = make(map[string]*apkInfo, 10000)
var newFilesMap = make(map[string]*apkInfo, 100)
var delFilesMap = make(map[string]*apkInfo, 100)
var biggerMap = make(map[string]*apkInfo, 100)
var smallerMap = make(map[string]*apkInfo, 100)

func hanleOldVersionFile(path string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	if f.IsDir() {
		return nil
	}
	timeDot := strings.Index(f.ModTime().String(), "+")
	if timeDot < 0 {
		timeDot = len(f.ModTime().String())
	}
	ai := &apkInfo{
		fileName:   f.Name(),
		modifyTime: f.ModTime().String()[:timeDot-1],
		size:       int64(f.Size()),
	}
	lowVersionMap[ai.fileName] = ai
	//println(path)
	return nil
}
func hanleNewVersionFile(path string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	if f.IsDir() {
		return nil
	}
	ai := &apkInfo{
		fileName:   f.Name(),
		modifyTime: f.ModTime().String()[:len(f.ModTime().String())-10],
		size:       int64(f.Size()),
	}
	highVersionMap[ai.fileName] = ai
	//println(path)
	return nil
}

func getFilelist(path string, f filepath.WalkFunc) {
	err := filepath.Walk(path, f)
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}
func CompareVersion() {
	for _, v := range lowVersionMap {
		value, ok := highVersionMap[v.fileName]
		if ok {
			//bigger than 1k
			if value.size-v.size > 0 {
				value.diffSize = value.size - v.size
				biggerMap[v.fileName] = value
			}
			if value.size-v.size < 0 {
				value.diffSize = value.size - v.size
				smallerMap[v.fileName] = value
			}
		} else {
			delFilesMap[v.fileName] = v
		}
	}
	for _, v := range highVersionMap {
		_, ok := lowVersionMap[v.fileName]
		if !ok {
			newFilesMap[v.fileName] = v
		}
	}
}
func ShowResult() {
	const dir = "data/"
	os.Mkdir(dir, 0777) //创建一个目录
	f, err := os.Create(dir + "output.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	var totalIncreased int64
	var totalDecreased int64
	fmt.Fprintf(f, "版本比较结果:%s,%s\n", path1, path2)
	f.WriteString("增大的文件如下:\n")
	for _, v := range biggerMap {
		if float64(v.diffSize/1024.0) > 10.0 {
			fmt.Fprintf(f, "%s,时间:%s,原始大小:%g Kb,比旧版本增长:%g Kb\n", v.fileName, v.modifyTime, float64(v.size/1024.0), float64(v.diffSize/1024.0))
		}
		totalIncreased += v.diffSize
	}
	f.WriteString("减小的文件如下:\n")
	for _, v := range smallerMap {
		if float64(v.diffSize/1024.0) > 5.0 {
			fmt.Fprintf(f, "%s,时间:%s,原始大小:%g Kb,比旧版本减小:%g Kb\n", v.fileName, v.modifyTime, float64(v.size/1024.0), float64(v.diffSize/1024.0))
		}
		totalDecreased += v.diffSize
	}
	f.WriteString("增加的文件如下:\n")
	for _, v := range newFilesMap {
		fmt.Fprintf(f, "%s,时间:%s,原始大小:%g Kb,比旧版本增长:%g Kb\n", v.fileName, v.modifyTime, float64(v.size/1024.0), float64(v.size/1024.0))
		totalIncreased += v.size
	}
	f.WriteString("删除的文件如下:\n")
	for _, v := range delFilesMap {
		fmt.Fprintf(f, "%s,时间:%s,原始大小:%g Kb,比旧版本缩小:%g Kb\n", v.fileName, v.modifyTime, float64(v.size/1024.0), float64(v.size/1024.0))
		totalDecreased += v.size
	}
	var diff int64 = totalIncreased - totalDecreased
	ws := fmt.Sprintf("总计增长:%g Kb,删除文件:%g Kb,变化:%g Kb\n", float64(totalIncreased/(1024.0)), float64(totalDecreased/(1024.0)), float64(diff/1024.0))
	f.WriteString(ws)
}

var (
	path1, path2 string
)

func main() {
	flag.Parse()
	path1 = flag.Arg(0)
	path2 = flag.Arg(1)
	fmt.Println("current path1 args:", path1)
	fmt.Println("current path2 args:", path2)
	getFilelist(path1, hanleOldVersionFile)
	getFilelist(path2, hanleNewVersionFile)
	CompareVersion()
	ShowResult()
}
