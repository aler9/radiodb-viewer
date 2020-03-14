#!/bin/sh -e

chown -R user:user /data

su - user -c "exec /build/db" &

su - user -c "exec /build/router" &

wait
