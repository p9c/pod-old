#!/bin/bash
echo $1 $2 $3 $4
podctl -s 127.0.0.1:11548 setgenerate $1 $2
podctl -s 127.0.0.1:11348 setgenerate $1 $3
podctl -s 127.0.0.1:11448 setgenerate $1 $4
