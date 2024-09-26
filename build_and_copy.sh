#!/bin/bash

go build -buildmode=plugin -trimpath -o ./modules/physics.so
cp modules/physics.so ~/Personal/nakama/data

