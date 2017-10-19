package main

import (
	"ac-common-go/glog"
	"ac-common-go/version"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func main() {
	ver := flag.Bool("version", false, "current version 1.0")
	flag.Parse()
	if *ver {
		version.ShowVersion(false)
		os.Exit(0)
	}
	version.ShowVersion(true)
	runtime.GOMAXPROCS(runtime.NumCPU())

	//	glog.SetLogLevel("INFO")
	err := ReadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer glog.Flush()
	server := NewHttpServer(CrmConf.ServiceAddr, CrmConf.ServicePort)
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
