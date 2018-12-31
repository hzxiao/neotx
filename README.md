## neotx

for creating neo transaction

### Who to use 

1. you should write the file `arg.json`
```json
{
    "version": 1,
    "type":"InvocationTransaction",
    "input": [
        {
            "hash": "b80f65fc5c0cc9a24ae2d613770202aae95dfa598f6541f75987b747eb5ca830",
            "value": 1000,
            "n": 0
        }
    ],
    "invocation":{
        "contract":"c88acaae8a0362cdbdedddf0083c452a3a8bb7b8",
        "operation": "transfer",
        "params": [
            "(address)ARbjp1wPh5XJchZpSjqHzGVQnnpTxNR1x7",
            "(address)APxpKoFCfBk8RjkRdKwyUnsBntDRXLYAZc",
            "(integer)1000000"
        ]
    },
    "from":"ARbjp1wPh5XJchZpSjqHzGVQnnpTxNR1x7",
    "fromPriKey":"L4RmQvd6PVzBTgYLpYagknNjhZxsHBbJq4ky7Zd3vB7AguSM7gF1",
    "to":"ARbjp1wPh5XJchZpSjqHzGVQnnpTxNR1x7",
    "assetId":"602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7",
    "value": 10
}
```

2. exec follow command
```shell
./neotx -f arg.json
```