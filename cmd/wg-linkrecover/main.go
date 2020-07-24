package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jeremiejig/wg-linkrecover/internal/wglinkrecover"
)

func main() {
	ifName := flag.String("ifname", "", "Interface name to monitor.")

	flag.Parse()
	if *ifName == "" {
		flag.Usage()
	}

	sigchan := make(chan os.Signal, 1)
	app := wglinkrecover.NewApp(*ifName)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	go handleSigStop(app, sigchan)
	app.Start()
}

func handleSigStop(app *wglinkrecover.App, sig <-chan os.Signal) {
	_ = <-sig
	app.Stop()
}
