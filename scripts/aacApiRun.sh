#!/bin/bash


export CGO_ENABLED=1
export GOARCH=arm64
export CC=aarch64-linux-gnu-gcc

cd /home/dfxp/Desktop/AutomatedAccessControl/Repos/USF.DFX.ASM.API
/usr/local/go/bin/go run main.go
