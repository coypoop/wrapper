#!/bin/sh
date=$(TZ=UTC date -r $(cd ../src; /usr/pkg/bin/git show -s --format=%ct) +%Y%m%d%H%MZ)
../src/build.sh -j$(/sbin/sysctl -n hw.ncpuonline) -B $date -M $PWD/obj/ $@
