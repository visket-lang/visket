# Visket

![](https://img.shields.io/github/workflow/status/arata-nvm/visket/Go?style=for-the-badge)
![](https://img.shields.io/github/license/visket-lang/visket?style=for-the-badge)
![](https://img.shields.io/codecov/c/github/arata-nvm/visket?style=for-the-badge)

A compiled programming language

## Example
```
func main() {
  print(fib(41))
  return 0
}

func fib(n) {
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
- [x] functions
- [ ] modules
- [x] if / else / then
- [x] for
- [x] while

### Types
- [x] int
- [ ] string
- [ ] struct
- [ ] array
- [ ] map
- [ ] bool
- [ ] func

## Dependencies
- Clang == 9.x
- GNU Make

## Development

### Building from source
1. `git clone https://github.com/arata-nvm/visket && cd visket`
2. `make`

### Compiling a Visket program
1. `./bin/visket -O <filename>`
