<h1>Ichain </h1>

[![version](https://img.shields.io/github/tag/icheckteam/ichain.svg)](https://github.com/icheckteam/ichain/releases/latest)
[![docs](https://readthedocs.org/projects/ichain/badge/?version=latest)](https://ichain.readthedocs.io/en/latest/?badge=latest)
[![license](https://img.shields.io/github/license/icheckteam/ichain.svg)](https://github.com/icheckteam/ichain/blob/master/LICENSE)
[![](https://tokei.rs/b1/github/icheckteam/ichain?category=lines)](https://github.com/icheckteam/ichain)

Branch    | Tests | Coverage
----------|-------|---------
develop   | [![CircleCI](https://circleci.com/gh/icheckteam/ichain/tree/develop.svg?style=shield)](https://circleci.com/gh/icheckteam/ichain/tree/develop) | [![codecov](https://codecov.io/gh/icheckteam/ichain/branch/develop/graph/badge.svg)](https://codecov.io/gh/icheckteam/ichain)
master    | [![CircleCI](https://circleci.com/gh/icheckteam/ichain/tree/master.svg?style=shield)](https://circleci.com/gh/icheckteam/ichain/tree/master) | [![codecov](https://codecov.io/gh/icheckteam/ichain/branch/master/graph/badge.svg)](https://codecov.io/gh/icheckteam/ichain)

English | [Vietnamese](README_VN.md)

Welcome to Ichain source code library!

Ichain is a blockchain based on tendermint. Ichain makes deploying, multiple networks connection and run supply chain application easier.

NOTE: The code is alpha version, but is in the process of rapid development. The master code may be unstable, stable version can be downloaded in the release page.

#### Need help?
- Submit issues via Github ichain repository
- Join [Telegram](https://t.me/ichainproject) where developers team is able to help you 

#### Features
- Supports thousands of transactions per second
- Quick block generation time
- Supply chain traceability
- Deploying and management product
- Scalable smart contract 
- Multiple networks connection
- Identification of digital identity.

### Modules

1. [Identity](https://github.com/icheckteam/documentation/blob/master/Identity.md)
2. [Asset](https://github.com/icheckteam/documentation/blob/master/Asset.md)
3. Shipping
4. Insurance
5. Invoice

### Minimum requirements

Requirement|Notes
---|---
Go version | Go1.9 or higher

### Install 

To download pre-built binaries, see our [Release page](https://github.com/icheckteam/ichain/releases)

Clone the ichain repository into the appropriable $GOPATH/src/github.com/icheckteam

```
$ git clone github.com/icheckteam/ichain.git
```

or 

```
$ go get github.com/icheckteam/ichain.git
```

Build the source with make.

```
$ make
```

After building the source code successfully. You should see two executable programs:

- `ichaind`: The node command line program for node control 
- `ichaincli`: The command line program for client control

### Public test network and sync node deployment

1. Create account 
- Through command line program, create an account
```
./ichaincli keys add testaccount
Enter a passphrase for your key:
Re-enter password:
testaccount     283873F09FEBC7EC95BCFBD43B37CF0678B8232A
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

melody tunnel slice calm basket round retreat cry impulse tail tunnel awkward morning wash apple abandon
```
2. Start ichain 
- Through command line program, start node
```
./ichaind start
```

### Implements

Run `ichaincli --help` for details

### Ichain transfer sample
```
./ichaincli transfer --name testaccount --amount 100tomato --to 283873F09FEBC7EC95BCFBD43B37CF0678B8232A
```
### Contribute

If you want to contribute please send email  (dev@icheck.vn)
