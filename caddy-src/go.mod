module caddy

go 1.16

replace github.com/gzzchh/caddy-forward-redirect => /data/SourceCode/caddy-follow-redirect

require (
	github.com/caddyserver/caddy/v2 v2.4.6
	github.com/caddyserver/forwardproxy v0.0.0-20210607180642-7b288bc31772
	github.com/gzzchh/caddy-forward-redirect v0.0.0-00010101000000-000000000000
	github.com/sjtug/caddy2-filter v0.0.0-20220212031450-0a3366033940
)
