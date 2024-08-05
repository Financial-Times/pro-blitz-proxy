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
    -H "X-Blitz-Cache-Id=uuid1" \
    -d '{"hello": "world"}' \
    -X POST \
    http://blitz/path/uri
```

### Get cache if it matches hash

```sh
curl \
    -H "X-Blitz-Cache-Id=uuid1" \
    -H "X-Blitz-Cache-Hash=hash" \
    -d '{"hello": "world"}' \
    -X POST \
    http://blitz/path/uri
```

### Get cache if it matches hash; invalidate older cache

```sh
curl \
    -H "X-Blitz-Cache-Id=uuid1" \
    -H "X-Blitz-Cache-Hash=hash" \
    -H "X-Blitz-Cache-Bust=older" \
    -d '{"hello": "world"}' \
    -X POST \
    http://blitz/path/uri
```

### Invalidate cache for specific hash

```sh
curl \
    -H "X-Blitz-Cache-Id=uuid1" \
    -H "X-Blitz-Cache-Hash=hash" \
    -H "X-Blitz-Cache-Bust=yes" \
    -d '{"hello": "world"}' \
    -X POST \
    http://blitz/path/uri
```

### Invalidate cache for key

```sh
curl \
    -H "X-Blitz-Cache-Id=uuid1" \
    -H "X-Blitz-Cache-Bust=yes" \
    -d '{"hello": "world"}' \
    -X POST \
    http://blitz/path/uri
```
