# Chainwalkers Elrond Parser

This tool communicates with the central Elrond Gateway (an Elrond Proxy instance) to fetch and parse blocks & transactions.

Chainwalkers Framework Docs: [docs.flipsidecrypto.com](https://docs.flipsidecrypto.com)

## Getting Started

## Building

In order to build the Docker image, run the following command:

```shell
cd v2
bash build_docker_image.sh 
```

This will create the image `chainwalkers_elrond:latest`.

## Up and Running

This parser has two entry points, each of which writes a json response to stdout.

1. `get_height.sh`: returns the current height of the Blockchain (actually, of the Elrond Metachain).
2. `parse_blocks.sh [int, int, ...]`: returns an array of blocks (actually, the Elrond Hyperblocks, an abstract structure derived from the Metablock), along with their transactions (if any).

To get the height of the blockchain, do as follows:

```shell
> docker run chainwalkers_elrond:latest bash get_height.sh
{
    "height": 12345
}
```

To parse a list of blocks run:

```shell
> docker run chainwalkers_elrond:latest bash parse_blocks.sh 1 2 3
[
    {
        "nonce": 1,
        "hash": "abba....",
        ...,
        "transactions": [
            {
                "nonce": 1,
                "hash": "baab...",
                ...
            },
            ...
        ],
        ...
    },
    ...
]
```
