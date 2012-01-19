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

const uploadform = `
<!doctype html>
<html>
<head>
<title>goof - upload</title>
<style type="text/css">
body{
font-family:"Lucida Grande", "Lucida Sans Unicode", Verdana, Arial, Helvetica, sans-serif;
font-size:12px;
}
p, h1, form, button{border:0; margin:0; padding:0;}
.spacer{clear:both; height:1px;}
.myform{
margin:0 auto;
width:400px;
padding:14px;
}
#stylized{
border:solid 2px #aaa;
background:#e9e9e9;
}
#stylized h1 {
font-size:14px;
font-weight:bold;
margin-bottom:8px;
}
#stylized p{
font-size:11px;
color:#666666;
margin-bottom:20px;
border-bottom:solid 1px #b7ddf2;
padding-bottom:10px;
}
#stylized label{
display:block;
font-weight:bold;
text-align:right;
width:140px;
float:left;
}
#stylized .small{
color:#666666;
display:block;
font-size:11px;
font-weight:normal;
text-align:right;
width:140px;
}
#stylized input{
float:left;
font-size:12px;
padding:4px 2px;
border:solid 1px #aacfe4;
width:200px;
margin:2px 0 20px 10px;
}
#stylized button{
clear:both;
margin-left:60px;
width:300px;
height:31px;
}
.button {
display: inline-block;
zoom: 1; /* zoom and *display = ie7 hack for display:inline-block */
*display: inline;
vertical-align: baseline;
margin: 0 2px;
outline: none;
cursor: pointer;
text-align: center;
text-decoration: none;
font: 14px/100% Arial, Helvetica, sans-serif;
padding: .5em 2em .55em;
text-shadow: 0 1px 1px rgba(0,0,0,.3);
-webkit-border-radius: .5em;
-moz-border-radius: .5em;
border-radius: .5em;
-webkit-box-shadow: 0 1px 2px rgba(0,0,0,.2);
-moz-box-shadow: 0 1px 2px rgba(0,0,0,.2);
box-shadow: 0 1px 2px rgba(0,0,0,.2);
}
.button:hover {
text-decoration: none;
}
.button:active {
position: relative;
top: 1px;
}
.bigrounded {
-webkit-border-radius: 2em;
-moz-border-radius: 2em;
border-radius: 2em;
}
.medium {
font-size: 12px;
padding: .4em 1.5em .42em;
}
.small {
font-size: 11px;
padding: .2em 1em .275em;
}
.gray {
color: #e9e9e9;
border: solid 1px #555;
background: #6e6e6e;
background: -webkit-gradient(linear, left top, left bottom, from(#888), to(#575757));
background: -moz-linear-gradient(top,  #888,  #575757);
filter:  progid:DXImageTransform.Microsoft.gradient(startColorstr='#888888', endColorstr='#575757');
}
.gray:hover {
background: #616161;
background: -webkit-gradient(linear, left top, left bottom, from(#757575), to(#4b4b4b));
background: -moz-linear-gradient(top,  #757575,  #4b4b4b);
filter:  progid:DXImageTransform.Microsoft.gradient(startColorstr='#757575', endColorstr='#4b4b4b');
}
.gray:active {
color: #afafaf;
background: -webkit-gradient(linear, left top, left bottom, from(#575757), to(#888));
background: -moz-linear-gradient(top,  #575757,  #888);
filter:  progid:DXImageTransform.Microsoft.gradient(startColorstr='#575757', endColorstr='#888888');
}
</style>
</head>
<body>
<div id="stylized" class="myform">
<form action="/upload" method="POST" id="uploadform" enctype="multipart/form-data">
<h1>Upload</h1>
<p>Powered by goof</p>
<label>File:
<span class="small">Choose file to upload</span>
</label>
<input type="file" id="fileinput" multiple="true" name="file">
<button class="button gray" type="submit">Send in the binary chaos</button>
<div class="spacer"></div>
</form>
<a target="_blank" href="https://github.com/stone/goof">goof</a>
</div>
</body>
`

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
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.RawPath)

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
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.RawPath)
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
