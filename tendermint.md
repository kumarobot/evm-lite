# Installation
```
go get github.com/bear987978897/evm-lite
cd ~/go/src/github.com/bear987978897/evm-lite
make vendor
make install
```
Then the system can run `evml`.

# Configuration

The application writes data and reads configuration from the directory specified
by the --datadir flag. The default data directory is `$HOME/.evm-lite`. The directory 
structure must respect the following stucture:

```
host:~/.evm-lite$ tree
├── tendermint
│    ├── config
│    │   ├── config.toml
│    │   ├── genesis.json
│    │   ├── node_key.json
│    │   └── priv_validator_key.json
│    └── data
│        └── priv_validator_state.json
├── eth
│   ├── genesis.json
│   ├── keystore
│   │   └── UTC--2018-10-14T11-12-24.412349157Z--633139fa62d5c27f454259ba59fc34773bd19457
│   └── pwd.txt
└── evml.toml
```

The `tendermint` directory can be generated by:  
`tendermint init [--home <evm-lite root directory>/tendermint]`.  
([tendermint](https://github.com/tendermint/tendermint) has to be installed)

The Ethereum genesis file defines Ethereum accounts and is stripped of 
all the Ethereum POW stuff. This file is useful to predefine a set of 
accounts that own all the initial Ether at the inception of the network.

Then running `evml tendermint` to start a single node.

# Multiple nodes

1. Command `tendermint testnet` can generate multiple tendermint files.
2. Place the generated files as the format in evm-lite. (see configuration section)
3. Modify `persistent peer` parameter to each node ip address in `tendermint/config/config.toml`.
4. Run nodes.