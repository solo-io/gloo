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