package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/vsukhin/numbers/controllers"
	"github.com/vsukhin/numbers/logger"
)

const (
	// PortHTTP contains server port
	PortHTTP = 8080
	// ParameterNamePortHTTP contains parameter name
	ParameterNamePortHTTP = "port"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
	httpPort = flag.Int(ParameterNamePortHTTP, PortHTTP, "HTTP server port")
)

func main() {
	logger.Log.Println("Starting server work at ", time.Now())
	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		close(c)
		logger.Log.Println("Finishing server work at ", time.Now())
		os.Exit(1)
	}()

	objectController := controllers.NewObjectControllerImplementation()

	mux := http.NewServeMux()
	mux.HandleFunc("/numbers", objectController.CleverGet)
	err := http.ListenAndServe(":"+strconv.Itoa(*httpPort), mux)
	if err != nil {
		logger.Log.Fatalf("Can't listen http port %v having error %v", httpPort, err)
	}
}
