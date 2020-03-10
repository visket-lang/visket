import "../lib/std"

struct Person {
  name: string
}

fun greet(person: Person) {
  printf("Hi, %s\n".cstring(), person.name)
}

fun main() {
  var person: Person
  person.name = "George"
  person.greet()
}

