
import requests

from walkpy import errors


class ElrondApi:
    def __init__(self, url: str):
        self.url = url

    def get_chain_height(self, shard: int) -> int:
        status = self.get_chain_status(shard)
        nonce = status.get("erd_highest_final_nonce", 0)
        return int(nonce)

    def get_chain_status(self, shard: int) -> dict:
        url = f"{self.url}/network/status/{shard}"
        response = _do_get(url)
        response = _open_proxy_response_envelope(url, response)
        response = response.get("status", dict())
        return response

    def get_block(self, nonce: int, shard: int):
        url = f"{self.url}/blocks?nonce={nonce}&shard={shard}"
        response = _do_get(url)
        return next(iter(response))


def _do_get(url):
    try:
        response = requests.get(url)
        response.raise_for_status()
        parsed = response.json()
        return parsed
    except requests.HTTPError as err:
        error_data = _extract_error_from_response(err.response)
        raise errors.ApiRequestError(url, error_data)
    except requests.ConnectionError as err:
        raise errors.ApiRequestError(url, err)
    except Exception as err:
        raise errors.ApiRequestError(url, err)


def _open_proxy_response_envelope(url, response: dict):
    err = response.get("error")
    code = response.get("code")

    if not err and code == "successful":
        return response.get("data", dict())

    raise errors.ApiRequestError(url, f"code:{code}, error: {err}")


def _extract_error_from_response(response):
    try:
        return response.json()
    except Exception:
        return response.text
