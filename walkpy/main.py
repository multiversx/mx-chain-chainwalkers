import logging
import sys
from argparse import ArgumentParser, FileType
from typing import Any

from walkpy import constants, utils
from walkpy.api import ElrondGateway


def main():
    logging.basicConfig(level=logging.INFO)

    parser = ArgumentParser()
    subparsers = parser.add_subparsers()

    sub = subparsers.add_parser("get-height")
    sub.add_argument("--outfile", type=FileType("w"), default=sys.stdout)
    sub.set_defaults(func=get_height)

    sub = subparsers.add_parser("parse-blocks")
    sub.add_argument("--outfile", type=FileType("w"), default=sys.stdout)
    sub.add_argument("--blocks", type=int, nargs='+')
    sub.set_defaults(func=parse_blocks)

    args = parser.parse_args()

    if not hasattr(args, "func"):
        parser.print_help()
    else:
        args.func(args)


def parse_blocks(args: Any):
    gateway = ElrondGateway(constants.GATEWAY_URL)
    outfile = args.outfile
    nonces = args.blocks

    blocks = []

    for nonce in nonces:
        hyperblock = gateway.get_hyperblock(nonce)
        blocks.append(hyperblock)

    utils.dump_json(blocks, outfile)


def get_height(args: Any):
    gateway = ElrondGateway(constants.GATEWAY_URL)
    outfile = args.outfile

    height = gateway.get_chain_height(constants.METASHARD_ID)
    utils.dump_json({"height": height}, outfile)


if __name__ == "__main__":
    main()
