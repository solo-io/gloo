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
FROM nginx:1.13.9-alpine
RUN rm -rf /etc/nginx/conf.d
COPY conf /etc/nginx
COPY --from=builder /usr/src/app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
