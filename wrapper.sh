#!/bin/sh
date=$(TZ=UTC date -r $(cd ../src; git show -s --format=%ct) +%Y%m%dT%H%M%SZ)
../src/build.sh -j$(sysctl -n hw.ncpuonline) -B $date $@
