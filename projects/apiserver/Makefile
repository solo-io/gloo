# These targets are for testing purposes only

NAME ?= "recorded_query"

.PHONY: record
record:
	NO_AUTH=1 OAUTH_SERVER="NO_SERVER_YET" OAUTH_CLIENT="NO_CLIENT_YET" DEBUG=1 DEBUG_RECORD_QUERIES=1 go run cmd/main.go

.PHONY: create-go-file
create-go-file:
	curl http://localhost:8082/querylog/go >> ./test/queries/${NAME}.go
