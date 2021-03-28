socks-proxy
===========
Simple socks5 proxy server with authentication and subnet restrictions.

## Config

Configuration is done with env variables.

| Variable | Default | Description |
| -------- | ------- | ----------- |
| PROXY_PORT | 1080 | Listen port for proxy |
| PROXY_USER | (empty) | Proxy username |
| PROXY_PASSWORD | (empty) | Proxy password |
| PROXY_BLOCK_DEST_NETS | (empty) | Comma-separated, without spaces, list of restricted destination subnets within proxy<br> Example: `127.0.0.0/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16` |
