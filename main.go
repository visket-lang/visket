package main

import (
  "fmt"
  "log"
  "os"
  "strconv"
)

func main() {
  if len(os.Args) != 2 {
    fmt.Println("Usage: solitude <number>")
    os.Exit(1)
  }

  n, err := strconv.Atoi(os.Args[1])
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  fmt.Println("define i32 @main() nounwind {")
  fmt.Printf("  ret i32 %d\n", n)
  fmt.Println("}")

}