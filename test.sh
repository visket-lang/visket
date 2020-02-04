#!/bin/bash

OPT=$*
TARGET=bin/solitude

try() {
  expected="$1"
  input="$2"

  echo "$input" > tmp.sl
  cat tmp.sl
  $TARGET $OPT -o tmp tmp.sl > /dev/null
  if [ "$?" != "0" ]; then
    exit 1
  fi
  ./tmp
  actual=$?

  if [ "$actual" == "$expected" ]; then
    echo "=> $actual"
  else
    echo "=> $expected expected, but got $actual"
    exit 1
  fi
}

try 0 "func main(): int { return 0 }"
try 42 "func main(): int { return 42 }"

try 4 "func main(): int { return 2+2 }"
try 52 "func main(): int { return 83-31 }"

try 10 "func main(): int { return 10   -   0 }"
try 10 "func main(): int { return 10 + 10 + 10 - 20 }"
try 0 "func main(): int { return 10 - 10 + 10 - 10 }"
try 20 "func main(): int { return 20 - 0 + 0 - 0 }"

try 30 "func main(): int { return 10 + 10 * 2 }"
try 102 "func main(): int { return 10 * 10 + 2 }"
try 200 "func main(): int { return 10 * 10 * 2 }"

try 40 "func main(): int { return (10 + 10) * 2 }"
try 120 "func main(): int { return 10 * (10 + 2) }"

try 10 "func main(): int { return (10 + 10) / 2 }"
try 5 "func main(): int { return 60 / (10 + 2) }"

try 1 "func main(): int { return 9 % 2 }"
try 3 "func main(): int { return 1 + 5 % 3 }"

try 16 "func main(): int { return 2 << 3 }"
try 2 "func main(): int { return 16 >> 3 }"

try 10 "func main(): int { return 120 + -110 }"
try 0 "func main(): int { return -(-10 - (-10)) }"

try 1 "func main(): int { if 10 == 10 { return 1 } return 0}"
try 0 "func main(): int { if 10 == 9 { return 1 } return 0}"
try 1 "func main(): int { if 10 != 9 { return 1 } return 0}"
try 0 "func main(): int { if 10 != 10 { return 1 } return 0}"

try 1 "func main(): int { if 9 < 10 { return 1 } return 0}"
try 0 "func main(): int { if 10 < 10 { return 1 } return 0}"
try 1 "func main(): int { if 10 <= 10 { return 1 } return 0}"
try 0 "func main(): int { if 10 <= 9 { return 1 } return 0}"
try 1 "func main(): int { if 10 > 9 { return 1 } return 0}"
try 0 "func main(): int { if 10 > 10 { return 1 } return 0}"
try 1 "func main(): int { if 10 >= 10 { return 1 } return 0}"
try 0 "func main(): int { if 9 >= 10 { return 1 } return 0}"

try 1 "func main(): int { if 9.0 < 10.0 { return 1 } return 0}"
try 0 "func main(): int { if 10.0 < 10.0 { return 1 } return 0}"
try 1 "func main(): int { if 10.0 <= 10.0 { return 1 } return 0}"
try 0 "func main(): int { if 10.0 <= 9.0 { return 1 } return 0}"
try 1 "func main(): int { if 10.0 > 9.0 { return 1 } return 0}"
try 0 "func main(): int { if 10.0 > 10.0 { return 1 } return 0}"
try 1 "func main(): int { if 10.0 >= 10.0 { return 1 } return 0}"
try 0 "func main(): int { if 9.0 >= 10.0 { return 1 } return 0}"

try 10 \
"func main(): int {
  var a = 10
  return a
}"

try 5 \
"func main(): int {
  var a = 10
  return a - 5
}"

try 10 \
"func main(): int {
  var a = 5
  a = a + 5
  return a
}"

try 7 \
"func main(): int {
  var a = 5
  a += 2
  return a
}"

try 3 \
"func main(): int {
  var a = 5
  a -= 2
  return a
}"

try 10 \
"func main(): int {
  var a = 5
  a *= 2
  return a
}"

try 2 \
"func main(): int {
  var a = 5
  a /= 2
  return a
}"

try 1 \
"func main(): int {
  var a = 5
  a %= 2
  return a
}"

try 12 \
"func main(): int {
  var a = 3
  a <<= 2
  return a
}"

try 5 \
"func main(): int {
  var a = 20
  a >>= 2
  return a
}"

try 2 "func num(): int { return 2 }
func main(): int { return num() }"
try 4 "func add(n: int): int { return n + 2 }
func main(): int { return add(2) }"

try 6 "func add(a: int, b: int): int { return a + b }
func main(): int { return add(2, 4) }"

try 2 \
"func main(): int {
  if true {
    return 2
  } else {
    return 1
  }
  return 0
}"

try 1 \
"func main(): int {
  if false {
    return 2
  } else {
    return 1
  }
  return 0
}"

try 2 \
"func main(): int {
  if true {
    if false {
      return 3
    } else {
      return 2
    }
  } else {
    return 1
  }
  return 0
}"

try 0 \
"func main(): int {
  if true {

  }
  return 0
}"

try 1 \
"func main(): int {
  if true {
    return 1
  } else {

  }

  return 2
}"

try 0 \
"func main(): int {
  if true {

  } else {
    return 1
  }

  return 0
}"

try 1 \
"func main(n: int): int {
  if n == 1 {
    return n
  } else {

  }
  return n
}"

try 55 \
"func fib(n: int): int {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}

func main(): int {
  return fib(10)
}"

try 45 \
"func main(): int {
  var sum = 0
  for var i = 0; i < 10; i = i + 1 {
    sum = sum + i
  }
  return sum
}"


try 45 \
"func main(): int {
  var sum = 0
  for var i = 0; i <= 9; i = i + 1 {
    sum = sum + i
  }
  return sum
}"

try 10 \
"func main(): int {
// return 0
  return 10 // 0
}"

try 10 \
"func main(): int {
  for
    var i = 0; // i < 5;
    i < 10; // i = i + 3
    i = i + 1 // {
  {
  }
  return i
}"

try 55 \
"func main(): int {
  var sum = 0
  for i in 0..10 {
    sum += i
  }
  return sum
}"

try 0 \
"func main(): int {
  var a: [3]int
  return a[0]
}"

try 3 \
"func main(): int {
  var a: [3]int
  var b = 2
  a[0] = 1
  a[b - 1] = 2
  a[b * 1] = b + 1
  return a[b]
}"

try 20 \
"func main(): int {
  return test1(10)
}

func test1(a: int): int {
  return a * 2
}"

echo "all tests passed"