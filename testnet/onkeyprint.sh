#!/bin/bash
while true;
do
	echo -ne "\033[0;0f"; ./printblocks.sh 0
	read -n1 -s
done
