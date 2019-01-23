go install
podctl -s 127.0.0.1:11348 setgenerate true 1
podctl -s 127.0.0.1:11448 setgenerate true 1
podctl -s 127.0.0.1:11548 setgenerate true 1
podctl -s 127.0.0.1:11648 setgenerate true 1
podctl -s 127.0.0.1:11748 setgenerate true 1
podctl -s 127.0.0.1:11848 setgenerate true 1
podctl -s 127.0.0.1:11948 setgenerate true 1
podctl -s 127.0.0.1:12048 setgenerate true 1
podctl -s 127.0.0.1:12148 setgenerate true 1
go install
podctl -s 127.0.0.1:11348 stop
podctl -s 127.0.0.1:11448 stop
podctl -s 127.0.0.1:11548 stop
podctl -s 127.0.0.1:11648 stop
podctl -s 127.0.0.1:11748 stop
podctl -s 127.0.0.1:11848 stop
podctl -s 127.0.0.1:11948 stop
podctl -s 127.0.0.1:12048 stop
podctl -s 127.0.0.1:12148 stop
go install
podctl -s 127.0.0.1:11348 setgenerate false 1
podctl -s 127.0.0.1:11448 setgenerate false 1
podctl -s 127.0.0.1:11548 setgenerate false 1
podctl -s 127.0.0.1:11648 setgenerate false 1
podctl -s 127.0.0.1:11748 setgenerate false 1
podctl -s 127.0.0.1:11848 setgenerate false 1
podctl -s 127.0.0.1:11948 setgenerate false 1
podctl -s 127.0.0.1:12048 setgenerate false 1
podctl -s 127.0.0.1:12148 setgenerate false 1
for i in $(seq 0 $(podctl -s 127.0.0.1:11348 getblockcount)); do podctl -s 127.0.0.1:11348 getblock `podctl -s 127.0.0.1:11348 getblockhash $i`|grep version\"|cut -f2 -d':'|cut -f1 -d","; done
echo $(for i in $(seq 0 $(podctl -s 127.0.0.1:11348 getblockcount)); do podctl -s 127.0.0.1:11348 getblock `podctl -s 127.0.0.1:11348 getblockhash $i`|grep version\"|cut -f2 -d':'|cut -f1 -d","; done)|xargs -n1 echo|sort|uniq
for i in $(seq 0 $(podctl -s 127.0.0.1:11348 getblockcount)); do podctl -s 127.0.0.1:11348 getblock `podctl -s 127.0.0.1:11348 getblockhash $i`; done
for i in $(seq 0 $(podctl -s 127.0.0.1:11348 getblockcount)); do podctl -s 127.0.0.1:11348 getblock `podctl -s 127.0.0.1:11348 getblockhash $i`|grep pow_hash\"|cut -f2 -d':'|cut -f1 -d","; done
for i in $(seq 0 $(podctl -s 127.0.0.1:11348 getblockcount)); do podctl -s 127.0.0.1:11348 getblock `podctl -s 127.0.0.1:11348 getblockhash $i`|grep -e version\" -e pow_hash\" -e time\" -e bits\" -e height\"|cut -f2 -d':'|cut -f1 -d","|xargs echo; done
for i in $(seq 0 $(podctl -s 127.0.0.1:11348 getblockcount)); do podctl -s 127.0.0.1:11348 getblock `podctl -s 127.0.0.1:11348 getblockhash $i`|grep time\"|cut -f2 -d':'|cut -f1 -d","; done
podctl -s 127.0.0.1:11348 getinfo
podctl -s 127.0.0.1:11448 getinfo
podctl -s 127.0.0.1:11548 getinfo
podctl -s 127.0.0.1:11648 getinfo
podctl -s 127.0.0.1:11748 getinfo
podctl -s 127.0.0.1:11848 getinfo
podctl -s 127.0.0.1:11948 getinfo
podctl -s 127.0.0.1:12048 getinfo
podctl -s 127.0.0.1:12148 getinfo
clear;rm -rf node0/testnet; ./testnet.sh node0
clear;rm -rf node1/testnet; ./testnet.sh node1
clear;rm -rf node2/testnet; ./testnet.sh node2
clear;rm -rf node4/testnet; ./testnet.sh node4
clear;rm -rf node5/testnet; ./testnet.sh node5
clear;rm -rf node6/testnet; ./testnet.sh node6
clear;rm -rf node8/testnet; ./testnet.sh node8
clear;rm -rf mining0/testnet; ./testnet.sh mining0
clear;rm -rf mining1/testnet; ./testnet.sh mining1
