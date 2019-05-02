FROM golang:1.12

COPY docker-entrypoint.sh /usr/local/bin/

RUN apt-get update -y && apt-get install jq -y && ln -s usr/local/bin/docker-entrypoint.sh /entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
