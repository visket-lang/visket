#!/bin/bash

OPT=$*
TARGET=bin/solitude

try() {
  expected="$1"
  input="$2"

  echo "$input" | $TARGET $OPT -o tmp.ll
  if [ "$?" != "0" ]; then
    echo "build failed"
    exit 1
  fi
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

try 10 "func main() { return (10 + 10) / 2 }"
try 5 "func main() { return 60 / (10 + 2) }"

try 1 "func main() { return 9 % 2 }"
try 3 "func main() { return 1 + 5 % 3 }"

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

try 10 "
func main() {
  var a = 10
  return a
}
"

try 5 "
func main() {
  var a = 10
  return a - 5
}
"

try 10 "
func main() {
  var a = 5
  a = a + 5
  return a
}
"

try 7 "
func main() {
  var a = 5
  a += 2
  return a
}
"

try 3 "
func main() {
  var a = 5
  a -= 2
  return a
}
"

try 10 "
func main() {
  var a = 5
  a *= 2
  return a
}
"

try 2 "
func main() {
  var a = 5
  a /= 2
  return a
}
"

try 1 "
func main() {
  var a = 5
  a %= 2
  return a
}
"

try 2 "func num() { return 2 }
func main() { return num() }"
try 4 "func add(n) { return n + 2 }
func main() { return add(2) }"

try 6 "func add(a, b) { return a + b }
func main() { return add(2, 4) }"

try 2 "
func main() {
  if 1 {
    return 2
  } else {
    return 1
  }
  return 0
}
"

try 1 "
func main() {
  if 0 {
    return 2
  } else {
    return 1
  }
  return 0
}
"

try 2 "
func main() {
  if 1 {
    if 0 {
      return 3
    } else {
      return 2
    }
  } else {
    return 1
  }
  return 0
}
"

try 0 "
func main() {
  if 1 {

  }
  return 0
}
"

# should be an error
#try 0 "
#func main() {
#  if 1 {
#    return 0
#  }
#
#}
#"

try 1 "
func main() {
  if 1 {
    return 1
  } else {

  }

  return 2
}
"

try 0 "
func main() {
  if 1 {

  } else {
    return 1
  }

  return 0
}
"

try 1 "
func main(n) {
  if n == 1 {
    return n
  } else {

  }
  return n
}
"

try 55 "
func fib(n) {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}

func main() {
  return fib(10)
}
"

try 45 "
func main() {
  var sum = 0
  for var i = 0; i < 10; i = i + 1 {
    sum = sum + i
  }
  return sum
}
"


try 45 "
func main() {
  var sum = 0
  for var i = 0; i <= 9; i = i + 1 {
    sum = sum + i
  }
  return sum
}
"

try 10 "
func main() {
// return 0
  return 10 // 0
}
"

try 10 "
func main() {
  for
    var i = 0; // i < 5;
    i < 10; // i = i + 3
    i = i + 1 // {
  {
  }
  return i
}
"
