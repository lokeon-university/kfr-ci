FROM python:3

COPY docker-entrypoint.sh /usr/local/bin/

RUN apt update && apt install python3-pip && ln -s usr/local/bin/docker-entrypoint.sh /

ENTRYPOINT ["docker-entrypoint.sh"]

