#!/bin/bash

go clean -testcache
for D in $(ls -1d sd*); do
  go test "github.com/gaorx/stardust5/$D"
#   if [ "$D" = "sdfile" ]; then
#     go test "github.com/gaorx/stardust4/$D/sdfiletype"
#     go test "github.com/gaorx/stardust4/$D/sdhttpfile"
#   elif [ "$D" = "sdcache" ]; then
#     go test "github.com/gaorx/stardust4/$D/sdcacheristretto"
#   fi
done

