to hack on gloo:

prereqs:

- minikube running


go run projects/sqoop/cmd/main.go &
go run projects/apiserver/cmd/main.go &
docker run -i -p 1234:8080 --rm soloio/petstore-example:latest &

