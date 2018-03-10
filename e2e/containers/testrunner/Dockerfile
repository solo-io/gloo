FROM ubuntu
RUN apt update && apt install -y curl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl \
    &&  chmod +x ./kubectl \
    && mv ./kubectl /usr/local/bin/
COPY test-ingress.crt /
CMD ["/bin/sh", "-c", "echo 'STARTING SLEEP! Access me.' && /bin/sleep 36000"]