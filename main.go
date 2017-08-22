package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	httpPort = "82"
	listenIP = ""
)

func main() {

	server := NewHttpServer(listenIP, httpPort)
	server.Serve()
}

type Server struct {
	ListenHost string
	ListenPort string
	Logger     *log.Logger
}

func NewHttpServer(listenIP, httpPort string) *Server {

	res := Server{
		ListenHost: listenIP,
		ListenPort: httpPort,
		Logger:     log.New(os.Stdout, "server> ", log.Ltime|log.Lshortfile)}

	domainHandler := NewDomainHandler()
	http.HandleFunc("/", res.HandleIndex)
	http.HandleFunc("/listDomain", domainHandler.handleListDomain)

	return &res
}

func (s *Server) Serve() {

	listenString := s.ListenHost + ":" + s.ListenPort
	s.Logger.Println("Serving http://" + listenString)
	s.Logger.Fatal(http.ListenAndServe(listenString, nil))
}

func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {

	defer timeTrack(time.Now(), w, "HandleIndex")

	s.httpHeaders(w)
	io.WriteString(w, "hello, world<br/><br/>")
}

func (s *Server) httpHeaders(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
}

func timeTrack(start time.Time, w http.ResponseWriter, name string) {

	elapsed := time.Since(start)
	io.WriteString(w, "<footer>"+name+" generated in "+elapsed.String()+"</footer>")
}
