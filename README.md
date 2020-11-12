Playing with TCP proxies.

Implementations:

- backend-go: an HTTP backend to proxy traffic to
- naive-tcp-proxy-go: spawn go routines and copy bytes
- naive-tcp-proxy-rs: spawn futures and copy bytes
- tcp-proxy-envoy: default envoy TCP proxy

Each proxy will implement the following features:

- Limit the number of open connections
- Stats for current number of open downstream and upstream connections
- Access logs with connection duration, bytes sent, bytes received, client address, and upstream address
- Upstream connection timeouts

Load testing:

TODO
