func fib(a, b, n) {
  if n > 0 {
    n = n - 1
    return fib(b, a + b, n)
  }
  return a
}

func main() {
  print(fib(0, 1, 46))
  return 0
}
