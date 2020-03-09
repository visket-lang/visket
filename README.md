
<h1 align="center">
    Visket
</h1>

<p align="center">
    <img src="https://img.shields.io/github/workflow/status/arata-nvm/visket/Go?style=for-the-badge" alt="GitHub Actions">
    <img src="https://img.shields.io/github/license/visket-lang/visket?style=for-the-badge" alt="Licence MIT">
    <img src="https://img.shields.io/codecov/c/github/arata-nvm/visket?style=for-the-badge" alt="Coverage">
</p>

<p align="center">
    A compiled programming language
</p>

<br>

## Example
```
func main() {
  print(fib(41))
}

func fib(n: int): int {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}
```

More examples can be found [here](https://github.com/arata-nvm/visket/tree/master/examples).

## Try it on Docker
1. `git clone https://github.com/arata-nvm/visket && cd visket`
2. `make docker/run`

## Features

### Language Features
- [x] variables
- [x] constants
- [x] functions
- [x] comments
- [x] modules
- [x] import
- [x] if / else / then
- [x] for
- [x] while
- [ ] if expression

### Types
- [x] bool
- [x] int
- [x] float
- [x] string
- [x] struct
- [x] array
- [ ] map
- [ ] func
- [ ] tagged union

## Dependencies
- Clang == 9.x
- GNU Make

## Development

### Building from source
1. `git clone https://github.com/arata-nvm/visket && cd visket`
2. `make`

### Compiling a Visket program
1. `./bin/visket -O <filename>`
