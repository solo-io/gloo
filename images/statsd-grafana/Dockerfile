FROM kamon/grafana_graphite@sha256:fc7517f21a44ac22f6ddb9feb6707372013ebdfaf532bae9edbdc92e200e8ff7

RUN     cd /src                                                                                                      &&\
        extract() { cat /opt/grafana/conf/custom.ini | grep $1 | awk '{print $NF}'; }                                &&\
        wizzy set grafana url $(extract ";protocol")://$(extract ";domain"):$(extract ";http_port")/grafana

RUN     sed -i 's@;root_url = %(protocol)s://%(domain)s:%(http_port)s/@root_url = %(protocol)s://%(domain)s:%(http_port)s/grafana@g' /opt/grafana/conf/custom.ini

COPY config.js /src/statsd/config.js