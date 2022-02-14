package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/caddyserver/forwardproxy"
	_ "github.com/gzzchh/caddy-forward-redirect"
	_ "github.com/sjtug/caddy2-filter"
)

func main() {
	caddycmd.Main()
}
