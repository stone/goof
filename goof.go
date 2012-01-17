package main

/* goof - a woof like server written in go
A very simple http server that serves files from
the current directory with the possibility to
receive files via upload form. 

Originally based on SimpleHttpd by mpl (github)
Written by Fredrik Steen <stone@ppo2.se>

Changelog:

v0.1:
  * Added possibility to make goof a download only server.
  * Added some sensible logging.
  * Added a starting banner.
  * Added error handler when starting up.
*/

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	host       = flag.String("host", "0.0.0.0:8080", "listening port and hostname")
	noUpload   = flag.Bool("n", false, "only allow downloads")
	help       = flag.Bool("h", false, "show this help")
	rootdir, _ = os.Getwd()
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: goof [flags]")
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
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.RawPath)
	http.ServeFile(w, r, path.Join(rootdir, url))
}

// A restricted version that only serves specified file
//func myFileServer(w http.ResponseWriter, r *http.Request, url string, filename string) {
//	fmt.Println("url:", url)
//	fmt.Println("filename:", filename)
//	http.ServeFile(w, r, path.Join(rootdir, url))
//
//}

func uploadHandler(w http.ResponseWriter, r *http.Request, url string) {
	mr, err := r.MultipartReader()

	if err != nil {
		// We did not receive any file so we serve up the form instead
		// or maybe we should redirect to a form handler?
		uploadFormHandler(w, r, url)
		return
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			http.Error(w, "reading body: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fileName := part.FileName()
		if fileName == "" {
			continue
		}

		log.Printf("%s %s %s", r.RemoteAddr, r.Method, fileName)

		buf := bytes.NewBuffer(make([]byte, 0))
		_, err = io.Copy(buf, part)
		if err != nil {
			http.Error(w, "copying: "+err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := os.Create(path.Join(rootdir, fileName))
		if err != nil {
			http.Error(w, "opening file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = buf.WriteTo(f)
		if err != nil {
			http.Error(w, "writing: "+err.Error(), http.StatusInternalServerError)
			return
		}

		break
	}
	http.Redirect(w, r, "/", 303)
}

// Called from myFileServer to serve the upload form
func uploadFormHandler(w http.ResponseWriter, r *http.Request, url string) {
	//TODO: make below a constant
	contents := `
<html>
<head>
  <title>Upload file</title>
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
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.RawPath)
	fmt.Fprintf(w, contents)
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

	// if flag is set we don't register the uploadHandler
	if *noUpload == false {
		http.HandleFunc("/upload", makeHandler(uploadHandler))
	}
	http.Handle("/", makeHandler(myFileServer))

	log.Printf("Serving %s on http://%s/", rootdir, *host)

	err := http.ListenAndServe(*host, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
