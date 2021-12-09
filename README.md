# File Share Server
This is a simple node/vuejs app for sharing files over network.

WARN: Files/Uploads are not locked down so this should only be used on a locked network.

## Run
Build
```bash
docker build -t file-share-server .
```
```bash
docker-compose up
```

## Testing
### Get files/folders
```bash
curl localhost:3000/api/files/...
```
## Delete files/folder
```bash
curl -X DELETE localhost:3000/api/files/...
```
## Create folder
```bash
curl -X POST localhost:3000/api/files/...
```
## Upload files
```bash
curl -X POST -F 'file=@/path/to/file' localhost:3000/api/files/...
```
