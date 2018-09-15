# EVM-LITE
## Ethereum with interchangeable consensus.

We took the [Go-Ethereum](https://github.com/ethereum/go-ethereum) 
implementation (Geth) and stripped out the EVM and Trie components to create a 
modular version that can work with different consensus systems. 

## ARCHITECTURE

```
                +-------------------------------------------+
+----------+    |  +-------------+         +-------------+  |       
|          |    |  | Service     |         | State App   |  |
|  Client  <-----> |             | <------ |             |  |
|          |    |  | -API        |         | -EVM        |  |
+----------+    |  | -Keystore   |         | -Trie       |  |
                |  |             |         | -Database   |  |
                |  +-------------+         +-------------+  |
                |         |                       ^         |     
                |         v                       |         |
                |  +-------------------------------------+  |
                |  | Engine                              |  |
                |  |                                     |  | 
                |  |       +----------------------+      |  |
                |  |       | Consensus            |      |  |
                |  |       +----------------------+      |  |
                |  |                                     |  |
                |  +-------------------------------------+  |
                |                                           | 
                +-------------------------------------------+

```

## Consensus Implementations:

- **SOLO**: No Consensus. Transactions are relayed directly from Service to 
            State
- **BABBLE**: Inmemory Babble node.

more to come...

## USAGE

Each consensus has its own subcommand `evml [consensus]`

```
Ethereum with interchangeable consensus

Usage:
  evml [command]

Available Commands:
  babble      Run the evm-lite node with Babble consensus
  help        Help about any command
  solo        Run the evm-lite node with Solo consensus (no consensus)
  version     Show version info

Flags:
      --datadir string        Top-level directory for configuration and data (default "/home/martin/.evm-lite")
      --eth.api_addr string   Address of HTTP API service (default ":8080")
      --eth.cache int         Megabytes of memory allocated to internal caching (min 16MB / database forced) (default 128)
      --eth.db string         Eth database file (default "/home/martin/.evm-lite/eth/chaindata")
      --eth.genesis string    Location of genesis file (default "/home/martin/.evm-lite/eth/genesis.json")
      --eth.keystore string   Location of Ethereum account keys (default "/home/martin/.evm-lite/eth/keystore")
      --eth.pwd string        Password file to unlock accounts (default "/home/martin/.evm-lite/eth/pwd.txt")
  -h, --help                  help for evml
      --log_level string      debug, info, warn, error, fatal, panic (default "debug")

Use "evml [command] --help" for more information about a command.

```

Instead of using CLI flags, options can also be specified in a `config.toml` 
file in the `datadir`. 

ex (config.toml):
``` toml
[eth]
db = "/eth.db"
[babble]
node_addr="127.0.0.1:1337"
```

## DEPLOY

We provide a set of scripts to automate the deployment of testnets. This 
requires [terraform](https://www.terraform.io/) and 
[docker](https://www.docker.com/).

ex: 
``` bash
cd deploy
#configure and start a testnet of 4 evm-lite nodes with Babble consensus
make consensus=babble nodes=4
#configure and start a single evm-lite instance with Solo consensus 
make consensu=solo nodes=1 
#bring everything down
make stop 
```

Support for AWS coming soon...


## DEV

To add a new consensus system:

- implement the consensus interface (consensus/consensus.go)
- add a property to the the global configuration object (config/config.go)
- create the corresponding CLI subcommand in cmd/evml/commands/
- register that command to the root command
