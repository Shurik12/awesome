FROM debian
SHELL [ "/bin/bash", "-c" ]
WORKDIR /root
RUN \
    apt update -y && apt upgrade -y && \
    apt install procps -y && \
    apt install lsof -y && \
    rm -rf /var/cache/apt/*
COPY  awesome /root/awesome
COPY  publish /root/publish
CMD sleep 15 && ./awesome
