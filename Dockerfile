FROM golang:1.13.8-alpine

RUN set -x \
  && apk update \
  && apk add --update --no-cache vim git make musl-dev curl \
  && apk add --update --no-cache clang llvm9 llvm9-dev binutils gcc

ADD . /visket

WORKDIR /visket

