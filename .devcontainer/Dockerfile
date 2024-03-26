FROM ubuntu:22.04

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN apt-get update --no-install-recommends \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
    && apt-get clean && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*

RUN apt-get update --no-install-recommends \
    && apt-get install -y --no-install-recommends \
        tzdata \
    && apt-get clean && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/* \
    && ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo 'Asia/Tokyo' >/etc/timezone

RUN apt-get update --no-install-recommends \
    && apt-get install -y \
        build-essential curl zip \
    && apt-get clean && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*

ARG GO_VERSION=1.21.0
RUN curl "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" -L --output go.tar.gz \
    && tar xvf go.tar.gz -C /usr/local/ \
    && echo 'export PATH=$PATH:/usr/local/go/bin:/home/vscode/.local/bin' >> /etc/profile \
    && rm go.tar.gz

RUN groupadd --gid 1000 vscode \
    && useradd --uid 1000 --gid 1000 -m vscode

USER vscode
ENV HOME /home/vscode
