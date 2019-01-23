#!/bin/bash
PORT=$2
for i in $(seq $1 $(podctl -s 127.0.0.1:$PORT getblockcount))
do podctl -s 127.0.0.1:$PORT getblock `podctl -s 127.0.0.1:$PORT getblockhash $i`|grep -e pow_hash\" -e time\" -e bits\" -e height\" -e \"pow_algo\"|cut -f2 -d':'|cut -f1 -d","|xargs echo
done
