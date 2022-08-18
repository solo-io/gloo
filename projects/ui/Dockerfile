FROM nginx:1.23.0-alpine

RUN apk -U upgrade

# Replace existing NGINX configuration
RUN rm -rf /etc/nginx/conf.d /etc/nginx/nginx.conf
COPY ./conf /etc/nginx

COPY ./build /usr/share/nginx/html

## add permissions for non-root user
RUN chown -R 10101:10101 /usr/share/nginx/html && \
        chmod -R 755 /usr/share/nginx/html && \
        chown -R 10101:10101 /var/cache/nginx && \
        chown -R 10101:10101 /var/log/nginx && \
        chown -R 10101:10101 /etc/nginx/conf.d
RUN touch /var/run/nginx.pid && \
        chown -R 10101:10101 /var/run/nginx.pid

USER 10101

EXPOSE 20202
CMD ["nginx", "-g", "daemon off;"]

