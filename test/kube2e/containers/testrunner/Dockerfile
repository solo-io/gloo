FROM ubuntu:22.04

# only get the things we say we want
RUN echo 'APT::Install-Suggests "0";' >> /etc/apt/apt.conf.d/00-keepclean
RUN echo 'APT::Install-Recommends "0";' >> /etc/apt/apt.conf.d/00-keepclean

# retry apt-get update and install to work around transient network issues
RUN echo 'Acquire::Retries "3";' >> /etc/apt/apt.conf.d/01-bad-sources
RUN echo 'Acquire::http::timeout "80";' >> /etc/apt/apt.conf.d/01-bad-sources
RUN echo 'Acquire::https::timeout "80";' >> /etc/apt/apt.conf.d/01-bad-sources

RUN apt update && apt install -y curl
COPY --from=lachlanevenson/k8s-kubectl:v1.10.3 /usr/local/bin/kubectl /usr/local/bin/kubectl

# Python
RUN apt-get install -y python; apt clean

COPY root.crt /

CMD ["/bin/sh", "-c", "echo 'STARTING SLEEP! Access me.' && /bin/sleep 36000"]