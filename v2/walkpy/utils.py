
import json
import sys
from typing import Any


def dump_json(data: Any, outfile: Any = None):
    if not outfile:
        outfile = sys.stdout

    json.dump(data, outfile, indent=4)
    outfile.write("\n")
