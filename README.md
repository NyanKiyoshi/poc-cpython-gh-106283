# PoC CPython GH-10XXXX

This repository contains a basic demonstration of a denial of service vulnerability
of `socket.create_connection()` against untrusted user-inputs of domain names.

[docker-compose.yaml](docker-compose.yaml) deploys a pre-configured DNS server (bind9)
with a `*.test` zone returning 256 IP addresses under the `A` record.

[poc.py](poc.py) invokes `socket.create_connection(("poc.test", 80), timeout=1)`
(default arguments) which simulates what would happen if a malicious user
would trick/ask a Python application to connect to a domain name returning
lots of IP addresses that always timeout when connecting to them.

### Usage

```
$ docker-compose up --build
$ docker-compose exec poc poc.py --domain poc.test --timeout 1
```

### Output

```
$ docker-compose exec poc poc.py --domain poc.test --timeout 1
2023-06-22 16:05:16,550 DEBUG: socket.connect: ('10.0.2.200', 80)
2023-06-22 16:05:17,553 DEBUG: socket.connect: ('10.0.2.84', 80)
[...]
2023-06-22 16:09:32,042 DEBUG: socket.connect: ('10.0.2.137', 80)
2023-06-22 16:09:33,043 ERROR: Socket timeout error!
2023-06-22 16:09:33,044 INFO: Took: 256 seconds
```

## Example non-vulnerable implementation

[example-golang.go](example-golang.go) shows the behavior of Golang versus
the Python's current, which acts an example on how the fix could be implemented
and would behave (RFC 8305).

```
$ docker-compose up --build
$ docker-compose exec go-example go-example --domain poc.test --timeout 10s
```

```
$ docker-compose exec go-example go-example --domain poc.test --timeout 10s
17:20:26.445149 Found 256 IPs to try out
17:20:26.445282 Attempting to connect to 10.0.2.21:80...
17:20:28.446960 Attempting to connect to 10.0.2.25:80...
17:20:30.449034 Attempting to connect to 10.0.2.67:80...
17:20:32.450936 Attempting to connect to 10.0.2.141:80...
17:20:34.452858 Attempting to connect to 10.0.2.13:80...
17:20:36.443316 Failed to get URL (in 10.004180651s): dial tcp 10.0.2.230:80: i/o timeout
```
