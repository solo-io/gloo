FROM kamon/grafana_graphite:8399049ce731
RUN     cd /src                                                                                                      &&\
        extract() { cat /opt/grafana/conf/custom.ini | grep $1 | awk '{print $NF}'; }                                &&\
        wizzy set grafana url $(extract ";protocol")://$(extract ";domain"):$(extract ";http_port")/grafana