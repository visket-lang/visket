#!/bin/bash

OPT=$*
TARGET=bin/visket

try() {
  expected="$1"
  input="$2"

  echo 'import "lib/testlib"' > tmp.sl
  echo "$input" >> tmp.sl
  cat tmp.sl
  $TARGET $OPT -o tmp tmp.sl > /dev/null
  if [ "$?" != "0" ]; then
    exit 1
  fi
  actual=`./tmp`

  if [ "$actual" == "$expected" ]; then
    echo "=> $actual"
  else
    echo "=> $expected expected, but got $actual"
    exit 1
  fi
}

try 0 "fun main() { print(0) }"
try 42 "fun main() { print(42) }"

try 4 "fun main() { print(2+2) }"
try 52 "fun main() { print(83-31) }"

try 10 "fun main() { print(10   -   0) }"
try 10 "fun main() { print(10 + 10 + 10 - 20) }"
try 0 "fun main() { print(10 - 10 + 10 - 10) }"
try 20 "fun main() { print(20 - 0 + 0 - 0) }"

try 30 "fun main() { print(10 + 10 * 2) }"
try 102 "fun main() { print(10 * 10 + 2) }"
try 200 "fun main() { print(10 * 10 * 2) }"

try 40 "fun main() { print((10 + 10) * 2) }"
try 120 "fun main() { print(10 * (10 + 2)) }"

try 10 "fun main() { print((10 + 10) / 2) }"
try 5 "fun main() { print(60 / (10 + 2)) }"

try 1 "fun main() { print(9 % 2) }"
try 3 "fun main() { print(1 + 5 % 3) }"

try 16 "fun main() { print(2 << 3) }"
try 2 "fun main() { print(16 >> 3) }"

try 10 "fun main() { print(120 + -110) }"
try 0 "fun main() { print(-(-10 - (-10))) }"

try 1 "fun main() { if 10 == 10 { print(1) } else { print(0) }}"
try 0 "fun main() { if 10 == 9 { print(1) } else { print(0) }}"
try 1 "fun main() { if 10 != 9 { print(1) } else { print(0) }}"
try 0 "fun main() { if 10 != 10 { print(1) } else { print(0) }}"

try 1 "fun main() { if 9 < 10 { print(1) } else { print(0) }}"
try 0 "fun main() { if 10 < 10 { print(1) } else { print(0) }}"
try 1 "fun main() { if 10 <= 10 { print(1) } else { print(0) }}"
try 0 "fun main() { if 10 <= 9 { print(1) } else { print(0) }}"
try 1 "fun main() { if 10 > 9 { print(1) } else { print(0) }}"
try 0 "fun main() { if 10 > 10 { print(1) } else { print(0) }}"
try 1 "fun main() { if 10 >= 10 { print(1) } else { print(0) }}"
try 0 "fun main() { if 9 >= 10 { print(1) } else { print(0) }}"

try 1 "fun main() { if 9.0 < 10.0 { print(1) } else { print(0) }}"
try 0 "fun main() { if 10.0 < 10.0 { print(1) } else { print(0) }}"
try 1 "fun main() { if 10.0 <= 10.0 { print(1) } else { print(0) }}"
try 0 "fun main() { if 10.0 <= 9.0 { print(1) } else { print(0) }}"
try 1 "fun main() { if 10.0 > 9.0 { print(1) } else { print(0) }}"
try 0 "fun main() { if 10.0 > 10.0 { print(1) } else { print(0) }}"
try 1 "fun main() { if 10.0 >= 10.0 { print(1) } else { print(0) }}"
try 0 "fun main() { if 9.0 >= 10.0 { print(1) } else { print(0) }}"

try 10 \
"fun main() {
  var a = 10
  print(a)
}"

try 5 \
"fun main() {
  var a = 10
  print(a - 5)
}"

try 10 \
"fun main() {
  var a = 5
  a = a + 5
  print(a)
}"

try 7 \
"fun main() {
  var a = 5
  a += 2
  print(a)
}"

try 3 \
"fun main() {
  var a = 5
  a -= 2
  print(a)
}"

try 10 \
"fun main() {
  var a = 5
  a *= 2
  print(a)
}"

try 2 \
"fun main() {
  var a = 5
  a /= 2
  print(a)
}"

try 1 \
"fun main() {
  var a = 5
  a %= 2
  print(a)
}"

try 12 \
"fun main() {
  var a = 3
  a <<= 2
  print(a)
}"

try 5 \
"fun main() {
  var a = 20
  a >>= 2
  print(a)
}"

try 2 "fun num(): int { return 2 }
fun main() { print(num()) }"
try 4 "fun add(n: int): int { return n + 2 }
fun main() { print(add(2)) }"

try 6 "fun add(a: int, b: int): int { return a + b }
fun main() { print(add(2, 4)) }"

try 2 \
"fun main() {
  if true {
    print(2)
    return;
  } else {
    print(1)
    return;
  }
  print(0)
}"

try 1 \
"fun main() {
  if false {
    print(2)
    return;
  } else {
    print(1)
    return;
  }
  print(0)
}"

try 2 \
"fun main() {
  if true {
    if false {
      print(3)
      return;
    } else {
      print(2)
      return;
    }
  } else {
    print(1)
    return;
  }
  print(0)
}"

try 0 \
"fun main() {
  if true {

  }
  print(0)
}"

try 1 \
"fun main() {
  if true {
    print(1)
    return;
  } else {

  }

  print(2)
}"

try 0 \
"fun main() {
  if true {

  } else {
    print(1)
    return;
  }

  print(0)
}"

try 0 \
"fun main() {
  var n: int
  if n == 1 {
    print(n)
    return;
  } else {

  }
  print(n)
}"

try 55 \
"fun fib(n: int): int {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}

fun main() {
  print(fib(10))
}"

try 10 \
"fun main() {
  var a = 5
  if a > 3 {
    print(10)
    return;
  } if a > 2 {
    print(5)
    return;
  }
  print(0)
}"

try 45 \
"fun main() {
  var sum = 0
  for var i = 0; i < 10; i = i + 1 {
    sum = sum + i
  }
  print(sum)
}"


try 45 \
"fun main() {
  var sum = 0
  for var i = 0; i <= 9; i = i + 1 {
    sum = sum + i
  }
  print(sum)
}"

try 10 \
"fun main() {
// print(0)
  print(10) // 0
}"

try 10 \
"fun main() {
  for
    var i = 0; // i < 5;
    i < 10; // i = i + 3
    i = i + 1 // {
  {
  }
  print(i)
}"

try 55 \
"fun main() {
  var sum = 0
  for i in 0..10 {
    sum += i
  }
  print(sum)
}"

try 0 \
"fun main() {
  var a: [3]int
  print(a[0])
}"

try 3 \
"fun main() {
  var a: [3]int
  var b = 2
  a[0] = 1
  a[b - 1] = 2
  a[b * 1] = b + 1
  print(a[b])
}"

try 20 \
"fun main() {
  print(test1(10))
}

fun test1(a: int): int {
  return a * 2
}"

try 10 \
"struct Foo {
  X: int
}
fun main() {
  var foo1 = new Foo
  foo1.X = 10
  var foo2 = new Foo
  foo2.X = 20
  print(foo1.X)
}"


try 20 \
"struct Foo {
  X: int
}
struct Bar {
  A: Foo
}
fun main() {
  var bar = new Bar
  bar.A = new Foo
  bar.A.X = 20
  print(bar.A.X)
}"

try 0 \
"struct Foo {
  X: int
}
fun main() {
  var foo1 = new Foo
  print(foo1.X)
}"

try 0 \
"struct Foo {
  X: int
}
struct Bar {
  A: Foo
}
fun main() {
  var bar = new Bar
  print(bar.A.X)
}"

try 10 \
"fun test(ref i: int) {
  i = 10
}
fun main() {
  var i = 0
  test(i)
  print(i)
}"

try 5 \
"fun test1(ref i: int) {
  test2(i)
}
fun test2(ref i: int) {
  i = 5
}
fun main() {
  var i = 0
  test1(i)
  print(i)
}"

try 10 \
"struct Foo {
  X: int
}
fun test(ref foo: Foo) {
  foo.X = 10
}
fun main() {
  var foo: Foo
  test(foo)
  print(foo.X)
}"

try 5 \
"var i = 5
fun main() { print(i) }"

try 0 \
"var i: int
fun main() { print(i) }"

try 10 \
"var i: int = 5
fun main() {
  i = 10
  print(i)
}"

try 9 \
"module Lib {
  fun test(): int { return 9 }
}
fun main() { print(Lib::test()) }"

try 4 \
"module Lib {
  fun test(): int { return 3 }
}
fun test(): int { return 1 }
fun main() { print(Lib::test() + test()) }"

try 1 \
"fun main() {
  var s = \"hoge\nfugapiyo\"
  if s[0] == 'h' {
    if s[4] == '\n' {
      print(1)
      return;
    }
  }
  print(0)
}"


echo "all tests passed"