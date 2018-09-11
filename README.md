# Faktur

A simple tool to timestamp data and relay them to somewhere.

Timestamping component.


Donation:

**BTC**: 1MYiMU3GfsgEP4EYHHonG9Fy6DkA1JC3B5

**ETH**: 0xc8f8371BDd6FB64388F0D65F43A0040926Ee38be

## Description

faktur-daemon is a connector to the Ethereum blockchain that ingest hashs and
create a RFC6962 Merkle Tree then generate receipts.

## Example

Through HTTP backend:

Save hashs:

`curl -v http://127.0.0.1:8090/save'?hash=0x42&hash=65&hash=65&hash=65&hash=65&hash=65'`

Retreive receipts:

`curl -s http://127.0.0.1:8090/getreceiptsbyroot'?target_hash=89c2b9c25146129a6424866925bce7d0ba6044e113219c3ffc062ac2651695e0' | python -m 'json.tool'`
```[
    {
        "ID": 5,
        "CreatedAt": "2018-09-11T02:22:41.096181Z",
        "UpdatedAt": "2018-09-11T02:22:41.096181Z",
        "DeletedAt": null,
        "Targethash": "0x89c2b9c25146129a6424866925bce7d0ba6044e113219c3ffc062ac2651695e0",
        "Proofs": "Ra90940cb21070fb27a3ec2e40ac2ac278cd4de182d4541a7f61ef4c91b8bb1eb",
        "MerkleRoot": "0xd4a2f8cd58f000b9cde00597e4775d460678f32af775ff3e4c02ed9446e1b1ad",
        "TransactionId": 3,
        "Transaction": {
            "ID": 3,
            "CreatedAt": "2018-09-11T02:22:41.071084Z",
            "UpdatedAt": "2018-09-11T02:22:53.077196Z",
            "DeletedAt": null,
            "MerkleRoot": "0xd4a2f8cd58f000b9cde00597e4775d460678f32af775ff3e4c02ed9446e1b1ad",
            "TransactionHash": "0xb55d93c3a79c8b264cf5d5032b0934b131a26dec1d2d63ec31aea84969b38d91",
            "Status": 3
        }
    }
]
```

## Receipt producer

 - [ ] Chainpoint2.1 based
 - [~] OpenTimestamps
 - [ ] BlockReceipt
 - [X] Custom

## Backends

 - [X] HTTP
 - [ ] Google Drive

## Delivrery system

 - [ ] Mail
 - [ ] Sentry
 - [ ] API (SMS, ...)
 - [X] Logging
