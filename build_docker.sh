#!/bin/bash

rm -rf modules
docker run --rm -w "/builder" --platform linux/amd64 -v "${PWD}:/builder" heroiclabs/nakama-pluginbuilder:3.22.0 build -buildmode=plugin -trimpath -o ./modules/physics.so
cp modules/physics.so ~/Personal/nakama/data
