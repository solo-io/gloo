# build environment
FROM node:10.13-alpine as builder

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

# If set (via --build-arg), will disable authentication in gloo-i
# This arg will be picked up during the `yarn install` step and injected into the application via DefinePlugin in Webpack configuration
ARG NO_AUTH

COPY ./gloo-i /usr/src/app/
RUN yarn install
RUN yarn build
RUN yarn buildglooui

# production environment
FROM nginx:1.13.9-alpine
RUN rm -rf /etc/nginx/conf.d
COPY conf /etc/nginx
COPY --from=builder /usr/src/app/packages/app/build /usr/share/nginx/html
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
