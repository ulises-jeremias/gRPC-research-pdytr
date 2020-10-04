ARG GOLANG_VERSION=latest

FROM golang:${GOLANG_VERSION} as base
# GOLANG_VERSION is specified again because the FROM directive resets ARGs
# (but their default value is retained if set previously)

ARG GOLANG_VERSION
ARG USER=test

# Needed for string substitution
SHELL ["/bin/bash", "-c"]

# Install base deps for development
RUN apt-get update \
    && apt-get -y install \
        sudo \
        neovim \
        apt-utils \
        locales \
    && rm -rf /tmp/* /var/tmp/*

RUN useradd -m ${USER} \
    && passwd -d ${USER} \
    && sed -i -e "s/Defaults    requiretty.*/ #Defaults    requiretty/g" /etc/sudoers \
    && echo "${USER} ALL=(ALL:ALL) NOPASSWD:ALL" > /etc/sudoers.d/${USER} \
    && usermod -a -G sudo ${USER} \
    && rm -rf /home/${USER}/.bashrc

# Install grpc
RUN go get -u google.golang.org/grpc \
    && go get -u github.com/golang/protobuf/protoc-gen-go

# Install protoc and zip system library
ARG PROTOC_VERSION=3.7.0
ENV PROTOC_FOLDER=protoc-${PROTOC_VERSION}-linux-x86_64
RUN apt-get update && apt-get install -y zip \
    && cd /tmp \
    && wget https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_FOLDER}.zip \
    && unzip ${PROTOC_FOLDER}.zip \
    && install -Dm755 bin/protoc /usr/bin/protoc

# Set correct locale
RUN echo "LC_ALL=en_US.UTF-8" >> /etc/environment \
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && echo "LANG=en_US.UTF-8" > /etc/locale.conf

RUN locale-gen en_US.UTF-8
ENV LC_CTYPE 'en_US.UTF-8'
ENV LANG C.UTF-8

ENV PATH=$PATH:$GOPATH/bin

COPY bashrc /etc/bash.bashrc
RUN chmod a+rwx /etc/bash.bashrc
