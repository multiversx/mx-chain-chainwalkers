# Chainwalkers Elrond Parser

This parser communicates with a Elrond database to extract and decode blocks & transactions.

Chainwalkers Framework Docs: [docs.flipsidecrypto.com](https://docs.flipsidecrypto.com)

## Getting Started

Open `./parsing/confing.toml` and set database URL, Username and Password

## Building

A docker image of the Elrond chainwalker can be built by running:

```shell
bash builder_docker_image.sh  (not done yet)
```

This will create an image with the tag `chainwalkers_elrond:latest`.

## Up and Running

This parser has two entry points, each of which writes a json response to stdout.

1. `get_height.sh`: returns the current height
2. `parse_blocks.sh [int, int, ...]`: returns an array of decoded blocks and each blocks transactions