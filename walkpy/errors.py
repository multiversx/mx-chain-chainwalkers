
from typing import Any


class KnownError(Exception):
    inner = None

    def __init__(self, message: str, inner: Any = None):
        super().__init__(message)
        self.inner = inner

    def get_pretty(self) -> str:
        if self.inner:
            return f"""{self}
... {self.inner}
"""
        return str(self)


class ApiRequestError(KnownError):
    def __init__(self, url: str, data: Any):
        super().__init__(f"Proxy request error for url [{url}]: {data}")
