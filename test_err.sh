#!/bin/sh

OPT=$*
TARGET=bin/visket

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
"fun main() {
  var a = 1
  a = 1.0
}"

try "tmp.sl:2 | undefined function 'notFound'" \
"fun main() {
  notFound()
}"

try "tmp.sl:2 | not enough arguments in call to 'test'" \
"fun main() {
  test()
}
fun test(a: int) {}"

try "tmp.sl:2 | too many arguments in call to 'test'" \
"fun main() {
  test(1, 1)
}
fun test(a: int) {}"

try "tmp.sl:2 | type mismatch 'float' and 'i32'" \
"fun main() {
  test(1.0)
}
fun test(a: int) {}"

try "tmp.sl:2 | unresolved variable 'a'" \
"fun main() {
  a
}"

try "tmp.sl:3 | cannot index 'i32'" \
"fun main() {
  var a = 1
  a[1]
}"

try "tmp.sl:2 | unexpected operator: float % float" \
"fun main() {
  1.0 % 1.0
}"

try "tmp.sl:2 | unexpected operator: i32.1" \
"fun main() {
  1 . 1
}"

try "tmp.sl:3 | unexpected operator: i32.A" \
"fun main() {
  var a = 1
  a.A
}"

try "tmp.sl:4 | unresolved member 'A'" \
"struct Foo { X: int }
fun main() {
  var foo = new Foo
  foo.A
}"

try "tmp.sl:2 | unknown type 'Hoge'" \
"fun main() {
  var a = new Hoge
}"

try "tmp.sl:2 | missing return at end of function" \
"fun test(): int {
}"

try "tmp.sl:2 | already declared function 'test'" \
"fun test() {}
fun test() {}"

try "tmp.sl:2 | type mismatch 'void' and 'i32'" \
"fun test() {
  return 1
}"

try "tmp.sl:2 | type mismatch 'i32' and 'float'" \
"fun test(): int {
  return 1.0
}"

try "tmp.sl:3 | already declared variable 'a'" \
"fun main() {
  var a = 1
  var a = 1
}"

try "tmp.sl:2 | type mismatch 'i32' and 'float'" \
"fun main() {
  var a: int = 1.0
}"

try "tmp.sl:1 | unknown type 'hoge'" \
"fun test(): hoge {
}"

try "tmp.sl:1 | illegal charactor '@'" \
"@"

try "tmp.sl:2 | type mismatch 'i32' and 'float'" \
"fun main() {
  for i in 0..1.0 {}
}"

try "tmp.sl:3 | a ref value must be an assignable variable" \
"fun test(ref i: int){}
fun main() {
  test(1)
}"

try "tmp.sl:4 | a ref value must be an assignable variable" \
"fun test(ref i: int){}
fun main() {
  val i = 1
  test(i)
}"

try "tmp.sl:1 | main func cannot have a return type" \
"fun main(): int {}"

try "tmp.sl:1 | main func cannot have parameters" \
"fun main(i: int) {}"

try "tmp.sl:3 | constant 'i' cannot be reassigned" \
"fun main() {
  val i = 10
  i = 1
}"

try "tmp.sl:2 | closing ' expected" \
"fun main() {
  var c = 'hoge'
}"

try "tmp.sl:2 | invalid escape sequence" \
"fun main() {
  var c = '\j'
}"

try "tmp.sl:4 | cannot load the member of incomplete structure: Foo.A" \
"struct Foo
fun main() {
  var foo: Foo
  var a = foo.A
}"

echo "all tests passed"