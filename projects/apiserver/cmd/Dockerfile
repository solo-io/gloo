FROM node:16.15-alpine
ARG GOARCH

RUN apk -U upgrade
COPY gloo-fed-apiserver-linux-${GOARCH} /usr/local/bin/gloo-fed-apiserver

COPY js /usr/local/bin/js
RUN npm install --prefix /usr/local/bin/js
USER 10101

ENTRYPOINT ["/usr/local/bin/gloo-fed-apiserver"]
