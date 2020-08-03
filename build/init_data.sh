#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

mkdir -p /data/go

for ARCH in ${ALL_ARCH}; do
    mkdir -p "/data/std/${ARCH}"
done

mkdir -p "${npm_config_cache}"

chown -R ${TARGET_UIDGID} /data
chmod -R a=rwX /data