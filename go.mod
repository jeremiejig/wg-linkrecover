module github.com/jeremiejig/wg-linkrecover

go 1.14

require (
	github.com/jsimonetti/rtnetlink v0.0.0-20200709124027-1aae10735293
	github.com/mdlayher/netlink v1.1.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sys v0.0.0-20200724161237-0e2f3a69832c
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20200609130330-bd2cb7843e1b
)

// replace github.com/jsimonetti/rtnetlink => github.com/jeremiejig/rtnetlink v0.0.0-20200724203954-4fe0743f2dee
