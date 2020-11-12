import logging
import sys
from argparse import ArgumentParser, FileType
from typing import Any

from walkpy import constants, errors, utils
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


def get_height(args: Any):
    gateway = ElrondGateway(constants.GATEWAY_URL)
    outfile = args.outfile

    height = gateway.get_chain_height(constants.METASHARD_ID)
    utils.dump_json({"height": height}, outfile)


def parse_blocks(args: Any):
    gateway = ElrondGateway(constants.GATEWAY_URL)
    outfile = args.outfile
    nonces = args.blocks

    blocks = []

    for nonce in nonces:
        try:
            hyperblock = gateway.get_hyperblock(nonce)
            post_process_hyperblock(hyperblock)
        except errors.ApiRequestError:
            hyperblock = {"nonce": nonce}
        blocks.append(hyperblock)

    utils.dump_json(blocks, outfile)


def post_process_hyperblock(block: dict):
    round = int(block.get("round", 0))
    block["timestamp"] = constants.ERD_START_TIME + constants.ERD_ROUND_TIME * round


if __name__ == "__main__":
    main()
