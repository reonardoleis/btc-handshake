# BTC Handshake
## How to test

- Download and run btcd (full node implementation written in Go), available (here)[https://github.com/btcsuite/btcd]
- Run this repository implementation with `make run`
- btcd should show `[INF] SYNC: New valid peer 127.0.0.1:PORT (inbound) (/rxonvrdo@challenge:1.0.0/)` on the stdout and close the connection afterwards