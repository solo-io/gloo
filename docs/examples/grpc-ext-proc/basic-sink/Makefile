.PHONY: docker-local server
all: server docker-local

server:
	@CGO_ENABLED=0 GOOS=linux go build -a --ldflags '-extldflags "-static"' \
		-tags "netgo" -installsuffix netgo \
		-o _output/server

docker-local: server
	docker build -t quay.io/solo-io/grpc-ext-proc-example .

# docker-push:
# 	docker push quay.io/solo-io/grpc-ext-proc-example

run:
	docker run -p 18080:18080 quay.io/solo-io/grpc-ext-proc-example

clean:
	@rm -fr _output/server


kload: 
	kind load docker-image quay.io/solo-io/grpc-ext-proc-example --name kind