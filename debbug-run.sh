#!/bin/bash
export CUR_WORK_DIR=/data/SourceCode/caddy-conf 
#./caddy-debug run --config ${CUR_WORK_DIR}/caddy/Caddyfile
/data/.buildTools/go/bin/dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./caddy-debug -- run --config ${CUR_WORK_DIR}/caddy/Caddyfile