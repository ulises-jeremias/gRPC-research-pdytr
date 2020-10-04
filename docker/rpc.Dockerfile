ARG UBUNTU_VERSION=latest

FROM ubuntu:${UBUNTU_VERSION} as base
# UBUNTU_VERSION is specified again because the FROM directive resets ARGs
# (but their default value is retained if set previously)

ARG UBUNTU_VERSION
ARG USER=test

# Needed for string substitution
SHELL ["/bin/bash", "-c"]

# Install base deps for development
RUN apt-get update \
    && apt-get -y install \
        build-essential \
        net-tools \
        sudo \
        neovim \
        apt-utils \
        locales \
        rpcbind \
        make \
    && rm -rf /tmp/* /var/tmp/*

ENV EDITOR nvim

RUN useradd -m ${USER} \
    && passwd -d ${USER} \
    && sed -i -e "s/Defaults    requiretty.*/ #Defaults    requiretty/g" /etc/sudoers \
    && echo "${USER} ALL=(ALL:ALL) NOPASSWD:ALL" > /etc/sudoers.d/${USER} \
    && usermod -a -G sudo ${USER} \
    && rm -rf /home/${USER}/.bashrc

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
