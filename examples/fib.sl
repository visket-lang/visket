func fib(n) {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}

func main() {
  print(fib(41))
  return 0
}
