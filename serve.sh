#!/bin/sh

./build.sh
caddy stop
caddy start --config ./Caddyfile