package wglinkrecover

import (
	"log"
	"os"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/jeremiejig/wg-linkrecover/internal/link"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// App is the entry point for wg-linkmonitor command
type App struct {
	interfaceName string

	wgClient  *wgctrl.Client
	stopped   int64
	linkState *link.Monitor
	// underlyingLinkState *link.UnderlyingMonitor
}

// NewApp return app to monitor the named interface
func NewApp(interfaceName string) *App {
	return &App{interfaceName: interfaceName}
}

// Start will start the app
func (app *App) Start() (status int) {
	var didPanic = true
	defer func() {
		if didPanic {
			e := recover()
			log.Printf("Panic %s:", e)
			debug.PrintStack()
			status = 2
		}
	}()
	var err error

	app.wgClient, err = wgctrl.New()
	if err != nil {
		panic(err)
	}
	wg, err := app.wgClient.Device(app.interfaceName)
	_ = wg
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("wireguard interface \"%s\" not found.", app.interfaceName)
			didPanic = false
			return 1
		}
		panic(err)
	}

	var d = 250 * time.Millisecond
	app.linkState, err = link.NewMonitor(app.interfaceName, d)
	if err != nil {
		panic(err)
	}
	// app.underlyingLinkState, err = link.NewUnderlyingMonitor(d*2, wg.Peers[0].Endpoint.IP, uint32(wg.FirewallMark))
	// if err != nil {
	// 	panic(err)
	// }
	app.main()
	didPanic = false
	return
}

// Close wrap app.Stop()
func (app *App) Close() {
	app.Stop()
}

// Stop stop the app
func (app *App) Stop() {
	atomic.StoreInt64(&app.stopped, 1)
}

func (app *App) main() {
	var portChanged bool
	var incremenPort int = 1
	// app.underlyingLinkState.Stop()
	defer app.linkState.Stop()
	defer app.wgClient.Close()
	for atomic.LoadInt64(&app.stopped) == 0 {
		select {
		case <-app.linkState.C:
			linkNotFound := atomic.LoadUint64(&app.linkState.LinkNotFoundTick)
			linkDowned := atomic.LoadUint64(&app.linkState.LinkDownedTick)
			if linkNotFound == 0 && linkDowned > 5 {
				log.Printf("%q link down !", app.interfaceName)
				if !portChanged {
					dev, err := app.wgClient.Device(app.interfaceName)
					if err != nil {
						// likely disappeared
						continue
					}
					conf := wgtypes.Config{
						ListenPort: new(int),
					}
					*conf.ListenPort = dev.ListenPort + incremenPort

					err = app.wgClient.ConfigureDevice(app.interfaceName, conf)
					if err == nil {
						portChanged = true
						incremenPort = 1
						log.Printf("%q new port: %d", app.interfaceName, *conf.ListenPort)
					} else {
						incremenPort++
					}
				}
			}
			if linkDowned == 0 {
				portChanged = false
			}
		}

	}
}
