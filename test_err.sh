#!/bin/bash

OPT=$*
TARGET=bin/solitude

try() {
  expected="$1"
  input="$2"

  echo "$input" > tmp.sl
  actual=`$TARGET $OPT -o tmp tmp.sl 2>&1 > /dev/null`
  # 'error:' を取り除く
  actual=${actual:16}
  if [ "$actual" == "$expected" ]; then
    echo "=> \"$actual\""
  else
    echo "=> \"$expected\" expected, but got \"$actual\""
    exit 1
  fi
}

#
try "tmp.sl:3 | type mismatch 'i32' and 'float'" \
"func main() {
  var a = 1
  a = 1.0
}"

try "tmp.sl:2 | undefined function 'notFound'" \
"func main() {
  notFound()
}"

try "tmp.sl:2 | not enough arguments in call to 'test'" \
"func main() {
  test()
}
func test(a: int) {}"

try "tmp.sl:2 | too many arguments in call to 'test'" \
"func main() {
  test(1, 1)
}
func test(a: int) {}"

try "tmp.sl:2 | type mismatch 'float' and 'i32'" \
"func main() {
  test(1.0)
}
func test(a: int) {}"

try "tmp.sl:3 | unresolved variable 'a'" \
"func main() {
  a
}"

try "tmp.sl:3 | cannot index 'i32'" \
"func main() {
  var a = 1
  a[1]
}"

try "tmp.sl:2 | unexpected operator: float % float" \
"func main() {
  1.0 % 1.0
}"

try "tmp.sl:2 | unexpected operator: i32.1" \
"func main() {
  1 . 1
}"

try "tmp.sl:3 | unexpected operator: i32.A" \
"func main() {
  var a = 1
  a.A
}"

try "tmp.sl:4 | unresolved member 'A'" \
"struct Foo { X: int }
func main() {
  var foo = new Foo
  foo.A
}"

try "tmp.sl:2 | unknown type 'Hoge'" \
"func main() {
  var a = new Hoge
}"

try "tmp.sl:2 | missing return at end of function" \
"func main(): int {
}"

try "tmp.sl:2 | already declared function 'main'" \
"func main() {}
func main() {}"

try "tmp.sl:2 | type mismatch 'void' and 'i32'" \
"func main() {
  return 1
}"

try "tmp.sl:2 | type mismatch 'i32' and 'float'" \
"func main(): int {
  return 1.0
}"

try "tmp.sl:3 | already declared variable 'a'" \
"func main() {
  var a = 1
  var a = 1
}"

try "tmp.sl:2 | type mismatch 'i32' and 'float'" \
"func main() {
  var a: int = 1.0
}"

try "tmp.sl:1 | unknown type 'hoge'" \
"func main(): hoge {
}"

try "tmp.sl:1 | Illegal charactor '@'" \
"@"

try "tmp.sl:2 | type mismatch 'i32' and 'float'" \
"func main() {
  for i in 0..1.0 {}
}"

echo "all tests passed"