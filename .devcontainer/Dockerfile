FROM ubuntu:20.04

ENV HOME /root
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
        build-essential curl zip fonts-noto-cjk \
        phantomjs \
    && apt-get clean && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*

ARG GO_VERSION=1.20.1
RUN curl "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" -L --output go.tar.gz \
    && tar xvf go.tar.gz -C /usr/local/ \
    && echo 'export PATH=$PATH:/usr/local/go/bin:/home/vscode/.local/bin' >> /etc/profile \
    && rm go.tar.gz

RUN groupadd --gid 1000 vscode \
    && useradd --uid 1000 --gid 1000 -m vscode

USER vscode
ENV HOME /home/vscode

RUN . /etc/profile \
    && go install -v github.com/ramya-rao-a/go-outline@latest \
    && go install -v github.com/cweill/gotests/gotests@latest \
    && go install -v github.com/fatih/gomodifytags@latest \
    && go install -v github.com/josharian/impl@latest \
    && go install -v github.com/haya14busa/goplay/cmd/goplay@latest \
    && go install -v github.com/go-delve/delve/cmd/dlv@latest \
    && go install -v honnef.co/go/tools/cmd/staticcheck@latest \
    && go install -v golang.org/x/tools/gopls@latest \
    && go install -v golang.org/x/tools/cmd/goimports@latest \
    && go install -v github.com/stamblerre/gocode@latest
