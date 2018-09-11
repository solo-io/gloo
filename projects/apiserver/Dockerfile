
# build environment
# FROM node:9.6.1 as builder
FROM alpine:3.6 as builder

#ARG UI_DIR

#RUN echo $UI_DIR

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

# COPY $UI_DIR/package.json $UI_DIR/yarn.lock /usr/src/app/
#RUN echo "ipv6" >> /etc/modules
run apk add yarn --no-cache
run yarn --version

COPY ./gloo-i /usr/src/app/
RUN ls -laht /usr/src/app/
run yarn install --frozen-lockfile --no-cache --production --ignore-engines
run yarn build
RUN ls -laht /usr/src/app/

# production environment
# FROM nginx:1.13.9-alpine
FROM alpine:3.6 as prod
RUN mkdir -p /bin/gloo-ui/build
ADD apiserver /bin/gloo-ui
COPY --from=builder /usr/src/app/build /bin/gloo-ui/build
EXPOSE 8082
RUN ls -laht /bin/gloo-ui/
RUN ls -laht /bin/gloo-ui/build/
#CMD ["/bin/gloo-ui/apiserver", "-d", "./build"]
CMD ["/bin/gloo-ui/apiserver", "-d", "/bin/gloo-ui/build"]
#CMD ["nginx", "-g", "daemon off;"]