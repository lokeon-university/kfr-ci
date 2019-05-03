FROM python:3
RUN apt update -y && apt install -y python3-pip
COPY docker-entrypoint.sh /usr/local/bin/
RUN ln -s /usr/local/bin/docker-entrypoint.sh /entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]