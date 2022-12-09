# two stage docker image -- the first stage clones envoy and
# the second builds the final image, copying necessary files from the first stage

# Stage 1: clone envoy
FROM alpine/git
WORKDIR /app
RUN git clone https://github.com/envoyproxy/envoy.git

# Stage 2: build the final image
FROM flask_service:python-3.10-slim-bullseye
# use the local responses.yaml file
ADD ./responses.yaml /etc/responses.yaml
# copy the service.py file from envoy
COPY --from=0 /app/envoy/examples/cache/service.py /code/service.py
ENV SERVICE_NAME=1