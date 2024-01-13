# Contracts

A repository containing smart contracts written in [Solidity](https://soliditylang.org/), compiled with [Hardhat](https://hardhat.org/), and tested using either Hardhat or the [Solo](https://wiki.iota.org/wasp-wasm/how-tos/solo/what-is-solo/) framework.

## Prerequisites

- [Go](https://go.dev/doc/install) `1.21 or above`
- [Node.js](https://nodejs.org) `16.9.0 or above` with [Corepack enabled](https://nodejs.org/api/corepack.html)

## Usage

To install dependencies and compile the contracts, run:

```
yarn
yarn compile
```

To run the tests, run:

```
yarn test:sandbox
```

or

```
yarn test:solo
```
