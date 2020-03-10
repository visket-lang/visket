import "../lib/std"

var memo: [100]int

fun main() {
  var result = fib(41)
  printf("%ld\n".cstring(), result)
}

fun fib(n: int): int {
  if memo[n] != 0 {
    return memo[n]
  }

  if n <= 1 {
    memo[n] = n
  } else {
    memo[n] = fib(n - 1) + fib(n - 2)
  }

   return memo[n]
}
