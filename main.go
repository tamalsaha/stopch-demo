package main


import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"time"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

func main() {
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := SetupSignalHandler()

	go f1(stopCh)
	go f2(stopCh)

	time.Sleep(10*time.Minute)
}

func f1(stopCh <-chan struct{}) {
	fmt.Println("f1")
	// go func() {
		<-stopCh
		fmt.Println("f1 >>>>>>>>>>>")
	//}()
}

func f2(stopCh <-chan struct{}) {
	fmt.Println("f2")
	//go func() {
		<-stopCh
		fmt.Println("f2 >>>>>>>>>>>")
	//}()
}
