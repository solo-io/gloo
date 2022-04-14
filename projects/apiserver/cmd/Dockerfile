FROM node:16.14.2-alpine

RUN apk -U upgrade
COPY gloo-fed-apiserver-linux-amd64 /usr/local/bin/gloo-fed-apiserver

COPY js /usr/local/bin/js
RUN npm install --prefix /usr/local/bin/js
USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed-apiserver"]
