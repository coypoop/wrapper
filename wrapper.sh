#!/bin/sh
../src/build.sh -j$(sysctl -n hw.ncpuonline) -B $(cd ../src; git show -s --format=%ct) $@
