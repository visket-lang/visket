#!/bin/bash

OPT=$*
TARGET=bin/solitude

try() {
  expected="$1"
  input="$2"

  echo "$input" | $TARGET $OPT -o tmp.ll
  lli tmp.ll
  actual=$?

  if [ "$actual" == "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

try 0 "func main() { return 0 }"
try 42 "func main() { return 42 }"

try 4 "func main() { return 2+2 }"
try 52 "func main() { return 83-31 }"

try 10 "func main() { return 10   -   0 }"
try 10 "func main() { return 10 + 10 + 10 - 20 }"
try 0 "func main() { return 10 - 10 + 10 - 10 }"
try 20 "func main() { return 20 - 0 + 0 - 0 }"

try 30 "func main() { return 10 + 10 * 2 }"
try 102 "func main() { return 10 * 10 + 2 }"
try 200 "func main() { return 10 * 10 * 2 }"

try 40 "func main() { return (10 + 10) * 2 }"
try 120 "func main() { return 10 * (10 + 2) }"

try 10 "func main() { return 120 + -110 }"
try 0 "func main() { return -(-10 - (-10)) }"

try 1 "func main() { return 10 == 10 }"
try 0 "func main() { return 10 == 9 }"
try 1 "func main() { return 10 != 9 }"
try 0 "func main() { return 10 != 10 }"

try 1 "func main() { return 9 < 10 }"
try 0 "func main() { return 10 < 10 }"
try 1 "func main() { return 10 <= 10 }"
try 0 "func main() { return 10 <= 9 }"
try 1 "func main() { return 10 > 9 }"
try 0 "func main() { return 10 > 10 }"
try 1 "func main() { return 10 >= 10 }"
try 0 "func main() { return 9 >= 10 }"
