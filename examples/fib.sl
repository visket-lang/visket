import "../lib/std"

fun main() {
  var result = fib(41)
  printi(result)
}

fun fib(n: int): int {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}
