FROM golang:alpine AS builder

RUN apk --no-cache add make
WORKDIR /app
COPY . .

FROM alpine
RUN ls
COPY ./_output/server /server

EXPOSE 18080
CMD ["./server"]
