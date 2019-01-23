#!/bin/bash
PORT=2008
for i in $(seq $1 $(podctl -s 127.0.0.1:$PORT getblockcount))
do podctl -s 127.0.0.1:$PORT getblock `podctl -s 127.0.0.1:$PORT getblockhash $i`|grep -e time\" |cut -f2 -d':'|cut -f1 -d",">>data
done
