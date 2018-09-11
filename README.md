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

`curl http://127.0.0.1:8090/getreceiptsbyroot'?target_hash=89c2b9c25146129a6424866925bce7d0ba6044e113219c3ffc062ac2651695e0' | python -m 'json.tool'`

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
