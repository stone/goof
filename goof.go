package main

/* goof - a woof like server written in go
A very simple http server that serves files from
the current directory with the possibility to
receive files via upload form.

Originally based on SimpleHttpd by mpl (github)
Written by Fredrik Steen <stone@ppo2.se>

Changelog:
v0.2:
  * Added option to stop serving after X requests
  * Added flag to only serve one file
  * Some cleanups.

v0.1:
  * Added possibility to make goof a download only server.
  * Added some sensible logging.
  * Added a starting banner.
  * Added error handler when starting up.

License:

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

const (
	onlyFileDefault      = "index.html"
	downloadCountDefault = 0
)


var (
	host          = flag.String("host", "0.0.0.0:8080", "listening port and hostname")
	noUpload      = flag.Bool("n", false, "only allow downloads")
	help          = flag.Bool("h", false, "show this help")
	onlyFile      = flag.String("f", onlyFileDefault, "restrict to one file")
	downloadCount = flag.Int("d", downloadCountDefault, "Max number of downloads")
	dcounter      = *downloadCount
	rootdir, _    = os.Getwd()
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: goof [flags]")
	flag.PrintDefaults()
	os.Exit(2)
}

func plog(r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.RequestURI())
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
	dcounter = dcounter + 1
	plog(r)

	// If downloads has reached max and downloadCount is not the default value
	if dcounter > *downloadCount && *downloadCount != downloadCountDefault {
		log.Fatal("Max downloads reached, quitting...")
	}

	// Serve only the file specified by user
	if *onlyFile != onlyFileDefault {
		http.ServeFile(w, r, path.Join(rootdir, *onlyFile))
		return
	}
	http.ServeFile(w, r, path.Join(rootdir, url))
}

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

		plog(r)

		buf := bytes.NewBuffer(make([]byte, 0))
		if _, err = io.Copy(buf, part); err != nil {
			http.Error(w, "copying: "+err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := os.Create(path.Join(rootdir, fileName))
		if err != nil {
			http.Error(w, "opening file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		if _, err = buf.WriteTo(f); err != nil {
			http.Error(w, "writing: "+err.Error(), http.StatusInternalServerError)
			return
		}

		break
	}
	http.Redirect(w, r, "/", 303)
}

// Called from myFileServer to serve the upload form
func uploadFormHandler(w http.ResponseWriter, r *http.Request, url string) {
	plog(r)
	fmt.Fprintf(w, uploadform)
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
	if *noUpload == false || *onlyFile != onlyFileDefault {
		http.HandleFunc("/upload", makeHandler(uploadHandler))
	}
	http.Handle("/", makeHandler(myFileServer))

	log.Printf("Serving %s on http://%s/", rootdir, *host)

	err := http.ListenAndServe(*host, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
