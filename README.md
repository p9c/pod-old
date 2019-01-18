# Parallelcoin Pod

Fully integrated all-in-one cli client, full node, wallet server, miner and GUI wallet for Parallelcoin

Pod is a multi-application with multiple submodules for different functions. It is self-configuring and configurations can be changed from the commandline as well as editing the json files directly, so the binary itself is the complete distribution for the suite.

It consists of 4 main modules:

1. pod/ctl - command line interface to send queries to a node or wallet and prints the results to the stdout
2. pod/node - full node for parallelcoin network, including optional indexes for address and transaction search, low latency miner controller
3. pod/wallet - wallet server that runs separately from the full node but depends on a full node RPC for much of its functionality, and includes a GUI front end
4. pod/shell - combined full node and wallet server with optional GUI

The shell is currently a simple wallet but will be expanded into a full application framework/shell.