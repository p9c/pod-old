#!/bin/bash
clear
pod --configfile=./node0/config --datadir=./node0 &
pod --configfile=./node1/config --datadir=./node1 &
pod --configfile=./node2/config --datadir=./node2 &
pod --configfile=./node3/config --datadir=./node3 &
pod --configfile=./node4/config --datadir=./node4 &
pod --configfile=./node5/config --datadir=./node5 &
pod --configfile=./node6/config --datadir=./node6 &
pod --configfile=./node7/config --datadir=./node7 &
pod --configfile=./node8/config --datadir=./node8 &
read -n1 -s
killall pod
