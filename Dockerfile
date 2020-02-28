FROM ubuntu:18.04 

RUN set -x \
  && apt-get update \
  && apt-get upgrade -y \
  && apt-get install software-properties-common -y \
  && add-apt-repository ppa:longsleep/golang-backports \
  && apt-get update \
  && apt-get install make clang-9 golang-1.13 -y \
  && ln -s /usr/bin/clang-9 /usr/bin/clang

ENV PATH /usr/lib/go-1.13/bin:$PATH

ADD . /solitude

WORKDIR /solitude

