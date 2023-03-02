package link

import (
	"sync/atomic"
	"time"

	"github.com/jsimonetti/rtnetlink"
	"github.com/mdlayher/netlink"
)

// Monitor hold a atomic variable which indicate how many tick the link is detected to be down.
type Monitor struct {
	// Must be read with atomic
	LinkDownedTick   uint64          // indicate how many tick the link did not receive packet (but sent packets)
	LinkNotFoundTick uint64          // indicate how many tick the Monitor associated link was not found (device does not exist)
	C                <-chan struct{} // When an update occured a struct{} is send to C

	c                 chan<- struct{}
	ifname            string
	ticker            *time.Ticker
	conn              *rtnetlink.Conn
	currentStats      rtnetlink.LinkStats64
	prevStats         rtnetlink.LinkStats64
	lastActivityStats rtnetlink.LinkStats64

	done chan struct{}
}

// NewMonitor return a new Monitor which will check the link statistics every d time.
func NewMonitor(ifname string, d time.Duration) (*Monitor, error) {
	m := &Monitor{
		ifname: ifname,
		done:   make(chan struct{}),
	}
	var err error
	var conf = netlink.Config{DisableNSLockThread: true}
	m.conn, err = rtnetlink.Dial(&conf)
	if err != nil {
		return nil, err
	}

	m.ticker = time.NewTicker(d)
	c := make(chan struct{}, 1)
	m.C = c
	m.c = c

	go m.main()
	return m, nil
}

func (m *Monitor) main() {
	defer m.conn.Close()
	m.updateStat() // Do one to get the tick zero
	for {
		select {
		case <-m.done:
			return
		case <-m.ticker.C:
			m.updateStat()
			if atomic.LoadUint64(&m.LinkNotFoundTick) == 0 {
				if m.currentStats.TXDropped != m.prevStats.TXDropped {
					atomic.AddUint64(&m.LinkDownedTick, 1)
				} else if m.currentStats.RXPackets != m.prevStats.RXPackets {
					// We received packet link is not down.
					atomic.StoreUint64(&m.LinkDownedTick, 0)
					// m.lastActivityStats = m.currentStats
				} else if m.currentStats.TXPackets != m.prevStats.TXPackets {
					atomic.AddUint64(&m.LinkDownedTick, 1)
				}
			} else {
				atomic.StoreUint64(&m.LinkDownedTick, 0)
			}
			select {
			case m.c <- struct{}{}:
			default:
			}
		}
	}
}

func (m *Monitor) updateStat() {
	links, err := m.conn.Link.List()
	if err != nil {
		return
	}
	for _, l := range links {
		if l.Attributes.Name == m.ifname {
			atomic.StoreUint64(&m.LinkNotFoundTick, 0)

			m.prevStats = m.currentStats
			m.currentStats = *l.Attributes.Stats64
			return
		}
	}
	atomic.AddUint64(&m.LinkNotFoundTick, 1)
}

// Stop will stop the monitor. After Stop, it cannot be started again.
func (m *Monitor) Stop() {
	m.ticker.Stop()
	close(m.done)
}
