FROM debian
SHELL [ "/bin/bash", "-c" ]
WORKDIR /root

COPY  awesome /root/
COPY  config.yml /root/
COPY  publish /root/publish

RUN \
    apt update -y && apt upgrade -y && \
    apt install procps -y && \
    apt install lsof -y && \
    rm -rf /var/cache/apt/*

CMD [ "tail", "-f", "/dev/null" ]
