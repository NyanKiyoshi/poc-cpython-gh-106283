#!/usr/bin/env python3
import argparse
import logging
import socket
import sys
import time

logger = logging.getLogger()
logger.setLevel(logging.DEBUG)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--domain", help="The domain to try connecting to", default="poc.test"
    )
    parser.add_argument(
        "--timeout", help="The timeout to set (integer or float)", type=float, default=1
    )
    return parser.parse_args()


def audit_hook(event, args):
    if event == "socket.connect":
        # Print socket address parameter whenever socket.connect() is invoked.
        logger.debug(f"{event}: {args[-1]}")


def main():
    sys.addaudithook(audit_hook)
    ns = parse_args()
    t0 = time.time()

    try:
        socket.create_connection((ns.domain, 80), timeout=ns.timeout)
    except socket.timeout:
        logger.error("Socket timeout error!")

    t1 = time.time()
    logger.info(f"Took: {int(t1 - t0)} seconds")


if __name__ == "__main__":
    logging.basicConfig(format="%(asctime)s %(levelname)s: %(message)s")
    main()
