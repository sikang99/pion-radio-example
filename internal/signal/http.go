package signal

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/sikang99/pion-radio-example/internal/common"
)

// createSessionKey is for the internal use to avoid duplaicate keys
func createSessionKey(key string) string {
	return "cojam" + key
}

// Middleware process context before handlers
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionid := r.Header.Get("x-cojam-session")
		sess, err := common.GetSession(sessionid)
		if err != nil {
			log.Fatalln(err)
		}

		// to control related workers
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ctx = context.WithValue(ctx, createSessionKey("session"), sess)
		nextRequest := r.WithContext(ctx)

		next(w, nextRequest)
	}
}

//
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
	log.Println(r.URL.String())
}

// PubHandler process publishers
func PubHandler(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value("session")
	if v == nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}
	sess := v.(common.SessionInfo)
	fmt.Fprintf(w, "PubHandler: "+string(sess.UserID))
	log.Println(r.URL.String())
}

// SubHandler process subscribers
func SubHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "SubHandler: ")
	log.Println(r.URL.String())
}

// MonHandler monitor the server
func MonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "MonHandler: ")
	log.Println(r.URL.String())
}

// HTTPSDPServer starts a HTTP Server that consumes SDPs
func HTTPSDPServer() (chan string, chan string, int) {
	port := flag.Int("port", 8080, "port of http server")
	rport := flag.Int("rport", 1234, "port of rtp data")
	dir := flag.String("dir", "static", "base directory of file server")
	//tout := flag.Int("time", 3, "timeout to serve in Second")
	flag.Parse()

	sdpInChan := make(chan string)
	sdpOutChan := make(chan string)

	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/pub", Middleware(PubHandler))
	http.HandleFunc("/sub", Middleware(SubHandler))
	http.HandleFunc("/mon", Middleware(MonHandler))

	http.HandleFunc("/sdp", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("/sdp connected from %s", r.Host)
		body, _ := ioutil.ReadAll(r.Body)
		// process request of sdp
		sdpInChan <- string(body)
		// send response of sdp
		fmt.Fprintf(w, <-sdpOutChan)
		log.Println("sent base64 SDP to client")
	})

	// http server for static files
	fs := http.FileServer(http.Dir(*dir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("\nWebRTC SFU example server is started")
	log.Printf("started http and file server on :%d and %s", *port, *dir)
	return sdpInChan, sdpOutChan, *rport
}
