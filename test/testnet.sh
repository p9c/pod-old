#!/bin/sh
# Parallelcoin isolated mainnet test harness
# Invoke with the following arguments:
# 
# testnet.sh datadir
pod --testnet --configfile=$1/config --datadir=$1 $2 
