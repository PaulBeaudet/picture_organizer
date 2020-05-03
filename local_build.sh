#!/bin/bash

# There is probably a beter way to do this, but it'll work on Mint Linux
go build organize.go
mkdir -p $HOME/bin
cp organize $HOME/bin
