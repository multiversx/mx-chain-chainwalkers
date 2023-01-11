# Chainwalkers MultiversX Parser

This parser communicates with a MultiversX database to extract and decode blocks & transactions.

Chainwalkers Framework Docs: [docs.flipsidecrypto.com](https://docs.flipsidecrypto.com)

## Getting Started

Open `./parsing/config.toml` and set database URL, Username and Password

## Building

A docker image of the MultiversX chainwalkers can be built by running:

```shell
cd v1
bash build_docker_image.sh 
```

This will create an image with the tag `chainwalkers_multiversx:latest`.

## Up and Running

This parser has two entry points, each of which writes a json response to stdout.

1. `get_height.sh`: returns the current height
2. `parse_blocks.sh [int, int, ...]`: returns an array of decoded blocks and each blocks transactions

To get the latest height of the blockchain run:

```shell
> docker run chainwalkers_multiversx:latest bash get_height.sh
{
    "height": 12345
}
```

To parse a list of blocks run:

```shell
> docker run chainwalkers_multiversx:latest bash parse_blocks.sh 1 2 3
[
    {
        "nonce": 1,
        "hash": "000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf",
        ...,
        "transactions": [
            {
                "nonce": 1,
                "hash": "000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf",
                ...
            },
            ...
        ],
        ...
    },
    ...
]
```
