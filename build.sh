#!/bin/bash
XCADDY_SKIP_CLEANUP=1 XCADDY_SKIP_BUILD=1 XCADDY_RACE_DETECTOR=1 XCADDY_DEBUG=1 xcaddy build \
    --output caddy-debug \
    --with github.com/caddyserver/forwardproxy@caddy2 \
    --with github.com/sjtug/caddy2-filter@master \
    --with github.com/gzzchh/caddy-forward-redirect=../caddy-follow-redirect