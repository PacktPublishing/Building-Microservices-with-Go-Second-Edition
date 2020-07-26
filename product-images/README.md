# Product Images

## Uploading 

Note: standard `-d` strips new lines

```
curl -vv localhost:9091/1/go.mod -X PUT --data-binary @$PWD/go.mod
```

## Downloading with compression

```
curl -v localhost:9091/1/go.mod --compressed -o file.tmp
```