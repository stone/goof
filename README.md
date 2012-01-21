#goof - share files through HTTP protocol

goof (Go Offer One File) is a very simple tool to send and receive files on your
local LAN. 

Features include:

 * Serve files just for a given number of times and then exit.
 * Share a directory and it's sub directories over HTTP.
 * Serve only one specified file.
 * Receive files via a simple html form to upload files via /upload
 * Disable upload functionality
 * More to come..


Usage:

    usage: goof [flags]
        -d=0: Max number of downloads
        -f="index.html": restrict to one file
        -h=false: show this help
        -host="0.0.0.0:8080": listening port and hostname
        -n=false: only allow downloads

Example:

    stone@ppo2:/tmp$ ./goof -host="0.0.0.0:8080"
    2012/01/17 16:15:10 Serving /tmp on http://0.0.0.0:8080/
    2012/01/17 16:15:12 127.0.0.1:38345 GET /
    2012/01/17 16:15:31 127.0.0.1:38347 GET /upload
    2012/01/17 16:15:44 127.0.0.1:38347 POST linux_3.2.0-7.13.tar.gz


To upload files with curl:

    stone@ppo2:/tmp$ curl -F file=@filetosend.iso  http://remotehost:8080/upload


Note 1: You need the [go][] runtime, <http://golang.org/>

Note 2: this is just a toy project in my adventures in the go language, it probably works
but not the cleanest code around ;) 

[go]:http://golang.org/  "The Go Programming language"
