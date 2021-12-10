package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/barklan/pg-vakt/pkg"
)

func handleSysSignals(data *pkg.Data) {
	SigChan := make(chan os.Signal, 1)

	signal.Notify(SigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-SigChan
	var sigID string
	switch sig {
	case syscall.SIGHUP:
		sigID = "SIGHUP"
	case syscall.SIGINT:
		sigID = "SIGINT"
	case syscall.SIGTERM:
		sigID = "SIGTERM"
	case syscall.SIGQUIT:
		sigID = "SIGQUIT"
	default:
		sigID = "UNKNOWN"
	}
	data.SendSync(fmt.Sprintf("I received %s. Exiting now!", sigID))
	time.Sleep(200 * time.Millisecond)
	data.B.Close()
	os.Exit(0)
}

func main() {
	log.Println("Starting...")

	data := pkg.InitData()
	defer data.SendSync("Deferred in main.")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func() {
			data.SendSync("Telebot poller exited.")
			wg.Done()
		}()

		pkg.RegisterHandlers(data.B, data)
		data.B.Start()
	}()

	go func() {
		handleSysSignals(data)
	}()

	go pkg.PeriodicDBBackups(data)
	go pkg.ContinuousDBBackups(data)

	data.Send("I am up!")
	wg.Wait()
}
