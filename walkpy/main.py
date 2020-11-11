import logging
import sys
from argparse import ArgumentParser, FileType
from typing import Any

from walkpy import constants, utils
from walkpy.api import ElrondApi


def main():
    logging.basicConfig(level=logging.INFO)

    parser = ArgumentParser()
    subparsers = parser.add_subparsers()

    sub = subparsers.add_parser("get-height")
    sub.add_argument("--url", default=constants.DEFAULT_URL)
    sub.add_argument("--outfile", type=FileType("w"), default=sys.stdout)
    sub.add_argument("--shard", type=int, default=constants.METASHARD_ID)
    sub.set_defaults(func=get_height)

    sub = subparsers.add_parser("parse-blocks")
    sub.add_argument("--url", default=constants.DEFAULT_URL)
    sub.add_argument("--outfile", type=FileType("w"), default=sys.stdout)
    sub.add_argument("--shard", type=int, default=constants.METASHARD_ID)
    sub.add_argument("--blocks", type=int, nargs='+')
    sub.add_argument("--with-validators", action="store_true")
    sub.set_defaults(func=parse_blocks)

    args = parser.parse_args()

    if not hasattr(args, "func"):
        parser.print_help()
    else:
        args.func(args)


def parse_blocks(args: Any):
    api = ElrondApi(args.url)
    outfile = args.outfile
    nonces = args.blocks
    shard = args.shard
    with_validators = args.with_validators

    blocks = []

    for nonce in nonces:
        block = api.get_block(nonce, shard)
        _post_process_block(block, with_validators)
        blocks.append(block)

    utils.dump_json(blocks, outfile)


def _post_process_block(block: dict, with_validators: bool):
    if "validators" in block and not with_validators:
        del block["validators"]


def get_height(args: Any):
    api = ElrondApi(args.url)
    outfile = args.outfile
    shard = args.shard

    height = api.get_chain_height(shard)
    utils.dump_json({"height": height}, outfile)


if __name__ == "__main__":
    main()
