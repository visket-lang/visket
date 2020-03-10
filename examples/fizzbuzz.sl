import "../lib/std"

fun main() {
  print("Plese input n: ")
  var n = inputi()

  for i in 1..n {
    if i % 15 == 0 {
      println("fizzbuzz")
    } else if i % 5 == 0 {
      println("buzz")
    } else if i % 3 == 0 {
      println("fizz")
    } else {
      printi(i)
    }
  }
}