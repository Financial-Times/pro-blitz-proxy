# Blitz Proxy 

Caching proxy with configurable storage

# Configure

```
# http port to be used by the proxy
PORT=4001
# http server to proxy
BACKEND=http://example.com:1234
```

## Storage

### File storage

If no configuration is provided, the default cache store is in files in the temp folder.

### Dynamodb

If DynamoDB `STORE_TYPE` is set and the table does not exist, it will be created automatically.

```
STORE_TYPE=DYNAMODB
AWS_ENDPOINT=http://localhost:28000
AWS_ACCESS_KEY_ID=local
AWS_SECRET_ACCESS_KEY=local
AWS_DYNAMODB_TABLENAME=test1
AWS_DYNAMODB_PROVISIONED_READ_CAPACITY=5
AWS_DYNAMODB_PROVISIONED_WRITE_CAPACITY=5
```

If you set AWS_DYNAMODB_PROVISIONED_READ_CAPACITY and AWS_DYNAMODB_PROVISIONED_WRITE_CAPACITY to 0, this will create DynamoDB table with 'per request' billing mode.

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
