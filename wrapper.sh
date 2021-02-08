#!/bin/sh
date=$(TZ=UTC date -r $(cd ../src; git show -s --format=%ct) +%Y%m%d%H%MZ)
../src/build.sh -j$(sysctl -n hw.ncpuonline) -B $date -M $PWD/obj/ $@
