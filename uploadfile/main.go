package main

import (
	//"bytes"
	"fmt"
	"io"
	//"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("Chunks"))))

	http.ListenAndServe(":1789", nil)
}

func upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	defer r.Body.Close()
	fileName := "DumpHitches" + time.Now().Format("20060102150405") + ".log"
	var reader io.Reader
	if r.Header.Get("Content-Type") == "text/plain" {
		f, err := os.OpenFile("Chunks/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		defer f.Close()
		reader = r.Body
		io.Copy(f, r.Body)
		fmt.Fprintln(w, "upload ok!")
		return
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		defer file.Close()
		fileName = handler.Filename
		reader = file
	}
	f, err := os.OpenFile("Chunks/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	defer f.Close()
	io.Copy(f, reader)
	fmt.Fprintln(w, "upload ok!")

}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(tpl))
}

const tpl = `<html>
<head>
<title>上传文件</title>
</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">
 <input type="file" name="uploadfile" />
 <input type="hidden" name="token" value="{...{.}...}"/>
 <input type="submit" value="upload" />
</form>
</body>
</html>`
