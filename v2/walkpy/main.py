import base64
import logging
import sys
from argparse import ArgumentParser, FileType
from typing import Any

from erdpy import errors, utils
from erdpy.accounts import Address
from erdpy.proxy import ElrondProxy

from walkpy import constants


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
    gateway = ElrondProxy(constants.GATEWAY_URL)
    outfile = args.outfile

    height = gateway.get_last_block_nonce(constants.METASHARD_ID)
    utils.dump_out_json({"height": height}, outfile)


def parse_blocks(args: Any):
    gateway = ElrondProxy(constants.GATEWAY_URL)
    outfile = args.outfile
    nonces = args.blocks

    blocks = []

    for nonce in nonces:
        try:
            hyperblock = gateway.get_hyperblock(nonce)
            post_process_hyperblock(hyperblock)
        except errors.ProxyRequestError:
            hyperblock = {"nonce": nonce}
        blocks.append(hyperblock)

    utils.dump_out_json(blocks, outfile)


def post_process_hyperblock(block: dict):
    # The timestamp can be computed deterministically, given the round and the genesis time
    round = int(block.get("round", 0))
    block["timestamp"] = constants.ERD_START_TIME + constants.ERD_ROUND_TIME * round

    # Hyperblocks V1 does not provide the transaction fees (nor "gasUsed"). Hyperblocks V2 will provide this information.
    # Until then, we can compute the fees deterministically (best-effort).
    # For Smart Contract transactions, the refunds will be incompletely provided by the API - until Hyperblocks V2 becomes available.
    # (e.g. intra-shard refunds are not visible etc.)

    for transaction in block["transactions"]:
        tx_type = transaction.get("type")
        tx_gas_limit = int(transaction.get("gasLimit", 0))
        tx_gas_price = int(transaction.get("gasPrice", 0))
        tx_receiver = Address(transaction.get("receiver"))
        tx_data = transaction.get("data", None)
        tx_data_decoded = ""
        tx_data_len = 0
        tx_gas_used = 0

        if tx_data:
            tx_data_decoded = base64.b64decode(tx_data)
            tx_data_len = len(tx_data_decoded)

        if tx_type == "normal":
            if tx_receiver.is_contract_address():
                tx_gas_used = tx_gas_limit
            else:
                tx_gas_used = constants.ERD_MIN_GAS_LIMIT + tx_data_len * constants.ERD_GAS_PER_DATA_BYTE

        fee = tx_gas_used * tx_gas_price
        transaction["fee"] = str(fee)


if __name__ == "__main__":
    main()
