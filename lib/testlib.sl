import "std"

fun print(i: int) {
  printf("%d\n".cstring(), i)
}

fun input(): int {
  var i: int
  scanf("%d".cstring(), i)
  return i
}