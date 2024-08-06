# Blitz Proxy 

Caching proxy with configurable storage

# Configure

Set backend address env variable:

```
BACKEND=http://example.com:1234
```

## Usage

### Get latest cache

```sh
curl \
    -H "X-Blitz-Cache-Id: uuid1" \
    -d '{"hello": "world"}' \
    -X POST \
    http://blitz/path/uri
```

### Make request without cache

```sh
curl \
    -H "X-Blitz-Cache-Id: uuid1" \
    -H "Cache-Control: no-cache" \
    -d '{"hello": "world"}' \
    -X POST \
    http://blitz/path/uri
```

### Invalidate cache for key

```sh
curl http://blitz/__internal/cache/bust/uuid1
```
