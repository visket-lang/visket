#!/bin/sh

OPT=$*
TARGET=bin/visket

try() {
  expected="$1"
  input="$2"

  echo 'import "lib/std"' > tmp.sl
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

try 0 "fun main() { printi(0) }"
try 42 "fun main() { printi(42) }"

try 4 "fun main() { printi(2+2) }"
try 52 "fun main() { printi(83-31) }"

try 10 "fun main() { printi(10   -   0) }"
try 10 "fun main() { printi(10 + 10 + 10 - 20) }"
try 0 "fun main() { printi(10 - 10 + 10 - 10) }"
try 20 "fun main() { printi(20 - 0 + 0 - 0) }"

try 30 "fun main() { printi(10 + 10 * 2) }"
try 102 "fun main() { printi(10 * 10 + 2) }"
try 200 "fun main() { printi(10 * 10 * 2) }"

try 40 "fun main() { printi((10 + 10) * 2) }"
try 120 "fun main() { printi(10 * (10 + 2)) }"

try 10 "fun main() { printi((10 + 10) / 2) }"
try 5 "fun main() { printi(60 / (10 + 2)) }"

try 1 "fun main() { printi(9 % 2) }"
try 3 "fun main() { printi(1 + 5 % 3) }"

try 16 "fun main() { printi(2 << 3) }"
try 2 "fun main() { printi(16 >> 3) }"

try 10 "fun main() { printi(120 + -110) }"
try 0 "fun main() { printi(-(-10 - (-10))) }"

try 1 "fun main() { if 10 == 10 { printi(1) } else { printi(0) }}"
try 0 "fun main() { if 10 == 9 { printi(1) } else { printi(0) }}"
try 1 "fun main() { if 10 != 9 { printi(1) } else { printi(0) }}"
try 0 "fun main() { if 10 != 10 { printi(1) } else { printi(0) }}"

try 1 "fun main() { if 9 < 10 { printi(1) } else { printi(0) }}"
try 0 "fun main() { if 10 < 10 { printi(1) } else { printi(0) }}"
try 1 "fun main() { if 10 <= 10 { printi(1) } else { printi(0) }}"
try 0 "fun main() { if 10 <= 9 { printi(1) } else { printi(0) }}"
try 1 "fun main() { if 10 > 9 { printi(1) } else { printi(0) }}"
try 0 "fun main() { if 10 > 10 { printi(1) } else { printi(0) }}"
try 1 "fun main() { if 10 >= 10 { printi(1) } else { printi(0) }}"
try 0 "fun main() { if 9 >= 10 { printi(1) } else { printi(0) }}"

try 1 "fun main() { if 9.0 < 10.0 { printi(1) } else { printi(0) }}"
try 0 "fun main() { if 10.0 < 10.0 { printi(1) } else { printi(0) }}"
try 1 "fun main() { if 10.0 <= 10.0 { printi(1) } else { printi(0) }}"
try 0 "fun main() { if 10.0 <= 9.0 { printi(1) } else { printi(0) }}"
try 1 "fun main() { if 10.0 > 9.0 { printi(1) } else { printi(0) }}"
try 0 "fun main() { if 10.0 > 10.0 { printi(1) } else { printi(0) }}"
try 1 "fun main() { if 10.0 >= 10.0 { printi(1) } else { printi(0) }}"
try 0 "fun main() { if 9.0 >= 10.0 { printi(1) } else { printi(0) }}"

try 10 \
"fun main() {
  var a = 10
  printi(a)
}"

try 5 \
"fun main() {
  var a = 10
  printi(a - 5)
}"

try 10 \
"fun main() {
  var a = 5
  a = a + 5
  printi(a)
}"

try 7 \
"fun main() {
  var a = 5
  a += 2
  printi(a)
}"

try 3 \
"fun main() {
  var a = 5
  a -= 2
  printi(a)
}"

try 10 \
"fun main() {
  var a = 5
  a *= 2
  printi(a)
}"

try 2 \
"fun main() {
  var a = 5
  a /= 2
  printi(a)
}"

try 1 \
"fun main() {
  var a = 5
  a %= 2
  printi(a)
}"

try 12 \
"fun main() {
  var a = 3
  a <<= 2
  printi(a)
}"

try 5 \
"fun main() {
  var a = 20
  a >>= 2
  printi(a)
}"

try 2 "fun num(): int { return 2 }
fun main() { printi(num()) }"
try 4 "fun add(n: int): int { return n + 2 }
fun main() { printi(add(2)) }"

try 6 "fun add(a: int, b: int): int { return a + b }
fun main() { printi(add(2, 4)) }"

try 2 \
"fun main() {
  if true {
    printi(2)
    return;
  } else {
    printi(1)
    return;
  }
  printi(0)
}"

try 1 \
"fun main() {
  if false {
    printi(2)
    return;
  } else {
    printi(1)
    return;
  }
  printi(0)
}"

try 2 \
"fun main() {
  if true {
    if false {
      printi(3)
      return;
    } else {
      printi(2)
      return;
    }
  } else {
    printi(1)
    return;
  }
  printi(0)
}"

try 0 \
"fun main() {
  if true {

  }
  printi(0)
}"

try 1 \
"fun main() {
  if true {
    printi(1)
    return;
  } else {

  }

  printi(2)
}"

try 0 \
"fun main() {
  if true {

  } else {
    printi(1)
    return;
  }

  printi(0)
}"

try 0 \
"fun main() {
  var n: int
  if n == 1 {
    printi(n)
    return;
  } else {

  }
  printi(n)
}"

try 55 \
"fun fib(n: int): int {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}

fun main() {
  printi(fib(10))
}"

try 10 \
"fun main() {
  var a = 5
  if a > 3 {
    printi(10)
    return;
  } if a > 2 {
    printi(5)
    return;
  }
  printi(0)
}"

try 45 \
"fun main() {
  var sum = 0
  for var i = 0; i < 10; i = i + 1 {
    sum = sum + i
  }
  printi(sum)
}"


try 45 \
"fun main() {
  var sum = 0
  for var i = 0; i <= 9; i = i + 1 {
    sum = sum + i
  }
  printi(sum)
}"

try 10 \
"fun main() {
// printi(0)
  printi(10) // 0
}"

try 10 \
"fun main() {
  for
    var i = 0; // i < 5;
    i < 10; // i = i + 3
    i = i + 1 // {
  {
  }
  printi(i)
}"

try 55 \
"fun main() {
  var sum = 0
  for i in 0..10 {
    sum += i
  }
  printi(sum)
}"

try 0 \
"fun main() {
  var a: [3]int
  printi(a[0])
}"

try 3 \
"fun main() {
  var a: [3]int
  var b = 2
  a[0] = 1
  a[b - 1] = 2
  a[b * 1] = b + 1
  printi(a[b])
}"

try 20 \
"fun main() {
  printi(test1(10))
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
  printi(foo1.X)
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
  printi(bar.A.X)
}"

try 0 \
"struct Foo {
  X: int
}
fun main() {
  var foo1 = new Foo
  printi(foo1.X)
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
  printi(bar.A.X)
}"

try 10 \
"fun test(ref i: int) {
  i = 10
}
fun main() {
  var i = 0
  test(i)
  printi(i)
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
  printi(i)
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
  printi(foo.X)
}"

try 5 \
"var i = 5
fun main() { printi(i) }"

try 0 \
"var i: int
fun main() { printi(i) }"

try 10 \
"var i: int = 5
fun main() {
  i = 10
  printi(i)
}"

try 9 \
"module Lib {
  fun test(): int { return 9 }
}
fun main() { printi(Lib::test()) }"

try 4 \
"module Lib {
  fun test(): int { return 3 }
}
fun test(): int { return 1 }
fun main() { printi(Lib::test() + test()) }"

try 1 \
"fun main() {
  var s = \"hoge\nfugapiyo\"
  if s[0] == 'h' {
    if s[4] == '\n' {
      printi(1)
      return;
    }
  }
  printi(0)
}"


echo "all tests passed"