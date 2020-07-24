package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jeremiejig/wireguard-monitorlink/internal/wglinkmonitor"
)

func main() {
	ifName := flag.String("iname", "", "Interface name to monitor.")

	flag.Parse()
	if *ifName == "" {
		flag.Usage()
	}

	sigchan := make(chan os.Signal, 1)
	app := wglinkmonitor.NewApp(*ifName)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	go handleSigStop(app, sigchan)
	app.Start()
}

func handleSigStop(app *wglinkmonitor.App, sig <-chan os.Signal) {
	_ = <-sig
	app.Stop()
}
