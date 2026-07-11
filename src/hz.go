package main
import (
  "fmt"
  "os"
)
func tukaikata() {
      fmt.Println("使い方をしれよ！！")
      fmt.Println("使い方は、hz [pack/unpack] ファイル名.その拡張子 変換するファイル名.hz(展開のときは、ファイル名.hz ファイル名.さっきまでの拡張子)")
      os.Exit(1)
}
func pack(input string,output string) {
  file,err := os.Create(output)
  if err != nil {
    fmt.Printf("圧縮ファイル作成エラーなり！！　詳細はこれなり！！ エラー内容:%s",err)
    return
    
  }
  inpfile,err := os.Open(input)
  if err != nil {
    fmt.Printf("圧縮元ファイル確認エラーなり！！　詳細はこれなり！！　エラー内容:%s",err)
    return
  }
  defer file.Close()
  defer inpfile.Close()
}
func main() {
  if len(os.Args) < 4 {
    tukaikata()
  }
  switch os.Args[1] {
    case "pack":
      pack(os.Args[2],os.Args[3])
    case "unpack":

    default:
      tukaikata()
  }
}
