ARG BASE=debian:stretch
FROM $BASE
 
LABEL maintainer="Mikhail Buslaev (buslaevnmh@yandex.ru)"
 
RUN apt-get update && \
    # Install basic utilities
    apt-get install --yes --allow-unauthenticated adduser vim sudo git curl unzip build-essential \
    # Install Compression libs
    zlib1g-dev libbz2-dev libsnappy-dev && \  
    # Cleanup
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Install GoLang
RUN curl -fsSL https://dl.google.com/go/go1.17.6.linux-amd64.tar.gz | tar xz \
    && chown -R root:root ./go && mv ./go /usr/local
ENV PATH="/usr/local/go/bin:${PATH}"

# Install DKV (Skipped for CI Pipelines)
ARG CI
RUN if [ -z "$CI" ] ; then git clone --depth=1 https://github.com/flipkart-incubator/dkv.git \
    && cd dkv && GOOS=linux GOARCH=amd64 make build \
    && mv ./bin /usr/lsocal/dkv && chown -R root:root /usr/local/dkv; fi

ENV PATH="/usr/local/exmo-trading:${PATH}"
ENV LD_LIBRARY_PATH="/usr/local/lib:$LD_LIBRARY_PATH"