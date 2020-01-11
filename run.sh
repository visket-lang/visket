#!/bin/bash

OPT=$*
TARGET=bin/solitude

$TARGET $OPT -o tmp.ll
if [ "$?" != "0" ]; then
  exit 1
fi

lli tmp.ll
