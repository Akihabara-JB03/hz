package main
import (
  "fmt"
  "os"
)
func main() {
  if len(os.Args) < 4 {
    fmt.Println("使い方は、hz [pack/unpack] ファイル名.その拡張子 変換するファイル名.hz(展開のときは、ファイル名.hz ファイル名.さっきまでの拡張子)")
    os.Exit(1)
  }
  switch os.Args[1] {
    case "pack":

    case "unpack":

    default:
      fmt.Println("使い方をしれよ！！")
      fmt.Println("使い方は、hz [pack/unpack] ファイル名.その拡張子 変換するファイル名.hz(展開のときは、ファイル名.hz ファイル名.さっきまでの拡張子)")
      os.Exit(1)
  }
}
