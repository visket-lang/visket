
<h1 align="center">
    <img src="https://raw.githubusercontent.com/visket-lang/design/master/logo.svg?sanitize=true">
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

## Docker上での実行方法

### コンパイラのビルド

以下のツールがインストールされていることを確認してください。

- Git
- GNU Make
- Docker

確認できれば、以下のコマンドを任意のターミナルにて実行します。

```
$ git clone https://github.com/visket-lang/visket && cd visket
$ make docker/run
```

処理が完了するとコンテナの中に入った状態になるので、そのまま以下のコマンドを実行します。

```
# make build
```

これでコンパイラがビルドされ、準備が完了しました。

### プログラムの実行

Visketのプログラムはコンテナ内で以下のコマンドを入力することにより実行できます。

```
# ./bin/visket -O -color <ファイル名>.sl && ./<ファイル名>
```

例として、`/visket/examples/hello_world.sl`を実行する際のコマンドは以下のようになります。
```
# ./bin/visket -O -color ./examples/hello_world.sl && ./hello_world
```

サンプルプログラムは`/visket/examples`ディレクトリ以下に保存されています。

まだαバージョンの段階であるため、バグを見つけた場合はIsuuesを立てていただけると私が助かります。

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
