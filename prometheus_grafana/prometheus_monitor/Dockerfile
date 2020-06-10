FROM centos7

WORKDIR /opt/prometheus_monitor

RUN mkdir -p /opt/prometheus_monitor/

ADD prometheus_monitor /opt/prometheus_monitor/

ENTRYPOINT ["/opt/prometheus_monitor/prometheus_monitor"]
