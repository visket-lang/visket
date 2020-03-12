import "../lib/std"

fun main() {
  var n = 40
  var result = fib(n)
  printf("%d\n".cstring(), result)
}

fun fib(n: int): int {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}

