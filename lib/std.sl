val false = (0 != 0)
val true = (0 == 0)

fun print(s: string) {
  printf(s.cstring())
}

fun println(s: string) {
  print(s)
  print("\n")
}

// TODO for test

fun inputi(): int {
  var i: int
  scanf("%d".cstring(), i)
  return i
}

fun inputf(): float {
  var f: float
  scanf("%f".cstring(), f)
  return f
}

fun inputd(): float64 {
  var f: float64
  scanf("%lf".cstring(), f)
  return f
}

fun printi(i: int) {
  printf("%d\n".cstring(), i)
}

fun printd(f: float64) {
  printf("%lf\n".cstring(), f)
}
