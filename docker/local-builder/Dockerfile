FROM docker.io/library/ubuntu:20.04

ARG USER_ID=1000
ARG GROUP_ID=1000

RUN useradd -u ${USER_ID} -g ${GROUP_ID} -mls /bin/bash builder

RUN  apt-get update && apt-get install -y curl \
  && curl -L https://go.dev/dl/go1.17.13.linux-$(dpkg --print-architecture).tar.gz -o /tmp/go.tar.gz \
  && rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz \
  && rm -rf /tmp/go.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
USER builder

RUN mkdir -p /home/builder/go && mkdir -p /home/builder/.cache

VOLUME [ "/home/builder/.cache" ]
VOLUME [ "/home/builder/go" ]
ENTRYPOINT [ "/bin/bash" ]

