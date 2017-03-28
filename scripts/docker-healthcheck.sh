#!/bin/sh
set -e

: ${CHAROND_PORT:=8080}
curl -f http://localhost:$((CHAROND_PORT+1))/health || exit 1