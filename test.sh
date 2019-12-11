#!/bin/bash

OPT=$*
TARGET=bin/solitude

try() {
  expected="$1"
  input="$2"

  echo "$input" | $TARGET $OPT > tmp.ll
  lli tmp.ll
  actual=$?

  if [ "$actual" == "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

try 0 "0"
try 42 "42"

try 4 "2+2"
try 52 "83-31"

try 10 "10   -   0"
try 10 "10 + 10 + 10 - 20"
try 0 "10 - 10 + 10 - 10"
try 20 "20 - 0 + 0 - 0"

try 30 "10 + 10 * 2"
try 102 "10 * 10 + 2"
try 200 "10 * 10 * 2"

try 40 "(10 + 10) * 2"
try 120 "10 * (10 + 2)"