package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

const uploadform = "upload.html"

var (
	host       = flag.String("host", "0.0.0.0:8080", "listening port and hostname")
	help       = flag.Bool("h", false, "show this help")
	rootdir, _ = os.Getwd()
)

func usage() {
	fmt.Fprintf(os.Stderr, "\t httpd \n")
	flag.PrintDefaults()
	os.Exit(2)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				return
			}
		}()
		title := r.URL.Path
		fn(w, r, title)
	}
}

// because getting a 404 when trying to use http.FileServer. beats me.
func myFileServer(w http.ResponseWriter, r *http.Request, url string) {
	http.ServeFile(w, r, path.Join(rootdir, url))
}

func uploadHandler(rw http.ResponseWriter, req *http.Request, url string) {
	mr, err := req.MultipartReader()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(rw, "reading body: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fileName := part.FileName()
		if fileName == "" {
			continue
		}
		println("fileName: " + fileName)
		buf := bytes.NewBuffer(make([]byte, 0))
		_, err = io.Copy(buf, part)
		if err != nil {
			http.Error(rw, "copying: "+err.Error(), http.StatusInternalServerError)
			return
		}
		f, err := os.Create(path.Join(rootdir, fileName))
		if err != nil {
			http.Error(rw, "opening file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		_, err = buf.WriteTo(f)
		if err != nil {
			http.Error(rw, "writing: "+err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
	http.ServeFile(rw, req, path.Join(rootdir, uploadform))
}

func createUploadForm() {
	contents := `
<html>
<head>
  <titleUpload file</title>
</head>
<body>
  <h1>Upload file</h1>

  <form action="/upload" method="POST" id="uploadform" enctype="multipart/form-data">
    <input type="file" id="fileinput" multiple="true" name="file">
    <input type="submit" id="filesubmit" value="Upload">
  </form>

</body>
</html>
`
	f, err := os.Create(path.Join(rootdir, uploadform))
	if err != nil {
		println("err creating uploadform")
		os.Exit(2)
	}
	defer f.Close()
	_, err = f.Write([]byte(contents))
	if err != nil {
		println("err writing uploadform")
		os.Exit(2)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if *help {
		usage()
	}

	nargs := flag.NArg()
	if nargs > 0 {
		usage()
	}

	createUploadForm()

	http.HandleFunc("/upload", makeHandler(uploadHandler))
	http.Handle("/", makeHandler(myFileServer))
	http.ListenAndServe(*host, nil)
}
