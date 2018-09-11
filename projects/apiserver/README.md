# Build
```bash
git clone https://github.com/solo-io/solo-kit
cd solo-kit
# pray
dep ensure -v 
go build -o apiserver projects/apiserver/cmd/main.go
./apiserver -h
```

## Run
```bash
./apiserver

## or with port

./apiserver -p 1234

```

By default runs on [http://localhost:8080](http://localhost:8080)

## Docker
- Get the latest version of the UI, build it (in a container), build the api server (locally), and copy built files into a docker container
```
export TAG=<versionNumber>; ./dockerScript
```
- push it to docker hub:
```
export TAG=<versionNumber>; docker push soloio/gloo-i:$TAG 
```



# TODO:
- auto generation of schema from protos
- dockerize / kuberize
- update when plugin api changes
