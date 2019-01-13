# neotx

for creating neo transaction

## What is neotx?

neotx is for creating neo transaction, include InvocationTransaction and ContactTransaction. then can send it to the neo network.

## Features

* Support InvocationTransaction and ContactTransaction
* Can send tx to neo node, include testnet, mainnet and prinet

## Installation

* if you want to install by source code
```go
go get -u github.com/hzxiao/neotx
```

* download with your browser from the [latest release](https://github.com/hzxiao/neotx/releases) page

## How to use 

1. you should write the file `arg.json`
```json
{
    "version": 1,
    "type": "InvocationTransaction",
    "input": [
        {
            "prevHash":"0x475fc8fb2b96b1ac6130d31cba20c7ecaf9b4d2d25c81314f2901f9c39479d31",
            "prevIndex": 0
        }
    ],
    "invocation": {
        "contract": "3805440cffa83a7d9509b1520e754a59a3ec579e",
        "operation": "transfer",
        "params": [
            "(address)AGHdThQFJs5kixWuXkgRsbNKz2LrDYDaQB",
            "(address)AWD7ju8oWGUMfpisa2ttFW6vEJYjdxSpZD",
            "(integer)300000000"
        ]
    },
    "from": "AGHdThQFJs5kixWuXkgRsbNKz2LrDYDaQB",
    "fromPriKey": "",
    "to": "AWD7ju8oWGUMfpisa2ttFW6vEJYjdxSpZD",
    "toPriKey":"",
    "assetId": "e13440dccae716e16fc01adb3c96169d2d08d16581cad0ced0b4e193c472eac1",
    "value": 100,
    "factor": 100000000,
    "doubleSign": true
}
```

2. example command
```shell
## create tx and print saw tx
neotx --arg arg.json

## create tx and send to testnet neo node
neotx --arg arg.json --send 

## crate tx and send to mainnet neo node
neotx --arg arg.json --net mainneet --send 

## create tx and send to private neo network
neotx --arg arg.json --send --net prinet --node http://localhost:20332
```

## Complete Command Line

```
Usage of neotx:
  -a, --arg string    a json format file is arg for creating a neo tx (default "arg.json")
  -h, --help          show usage of neotx command
  -n, --net string    give a network(testnet|mainnet|prinet) of neo (default "testnet")
  -d, --node string   give a node of neo if necessary. it is necessary when network is prinet
  -s, --send          send rpc request to neo node
  -v, --version       print version
```