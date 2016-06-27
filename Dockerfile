FROM scratch
MAINTAINER Kelsey Hightower <kelsey.hightower@gmail.com>
ADD mysql-healthz /mysql-healthz
ENTRYPOINT ["/mysql-healthz"]
