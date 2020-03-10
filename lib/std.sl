val false = (0 != 0)
val true = (0 == 0)

fun print(s: string) {
  printf(s.cstring())
}

fun println(s: string) {
  print(s)
  print("\n")
}

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
  scanf("%f".cstring(), f)
  return f
}

// TODO for test
fun printi(i: int) {
  printf("%d\n".cstring(), i)
}
