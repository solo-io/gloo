# build environment
FROM alpine:3.6 as builder

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

RUN apk add yarn --no-cache
RUN yarn --version

COPY ./gloo-i /usr/src/app/
RUN yarn install --frozen-lockfile --no-cache --production --ignore-engines
RUN yarn build

# production environment
FROM alpine:3.6 as prod
RUN mkdir -p /bin/gloo-ui/build
ADD apiserver /bin/gloo-ui
COPY --from=builder /usr/src/app/build /bin/gloo-ui/build
EXPOSE 8082
CMD ["/bin/gloo-ui/apiserver", "-d", "/bin/gloo-ui/build"]