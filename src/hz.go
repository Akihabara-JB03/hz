package main
import (
  "fmt"
  "os"
  "io"
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
  defer file.Close()
  inpfile,err := os.Open(input)
  
  
  
  if err != nil {
    fmt.Printf("圧縮元ファイル確認準備エラーなり！！　詳細はこれなり！！　エラー内容:%s",err)
    return
  }
  defer inpfile.Close()
  data, err := io.ReadAll(inpfile)
  if err != nil {
    fmt.Printf("圧縮元ファイル確認エラーなり！！詳細はこれなり！！　エラー内容:%s",err)
    return
  }
  count := 1
  s := string(data)
  for i := 0; i < len(s); i++{
    if i+1 < len(s) && s[i] == s[i+1] {
      // 同じ文字
      count++
    } else {
      if count >= 5 {
        fmt.Fprintf(file, "%c{%d}", s[i], count)
      } else {
        for j := 0; j < count; j++ {
          fmt.Fprintf(file, "%c", s[i])
        }
      }
      count = 1
    }
   
  }
   fmt.Println("圧縮完了なり！！")
}
func unpack(input string, output string) {
	file, err := os.Create(output)
	if err != nil {
		fmt.Printf("展開ファイル作成エラーなり！！ エラー内容:%s\n", err)
		return
	}
	defer file.Close()

	inpfile, err := os.Open(input)
	if err != nil {
		fmt.Printf("展開ファイル確認準備エラーなり！！ エラー内容:%s\n", err)
		return
	}
	defer inpfile.Close()

	data, err := io.ReadAll(inpfile)
	if err != nil {
		fmt.Printf("データ読み込みエラーなり！！ エラー内容:%s\n", err)
		return
	}

	s := string(data)

	for i := 0; i < len(s); i++ {
		// もし次の文字が '{' だったら、カッコの中の数字を読み取る
		if i+1 < len(s) && s[i+1] == '{' {
			currentChr := s[i] // 繰り返したい文字（例: 'A'）

			var count int
			var readLen int
			// カッコの中の数字（例: 300）と、その文字数（readLen）を同時に取得
			_, err := fmt.Sscanf(s[i+1:], "{%d}", &count)

			if err == nil {
				// 読み取れた回数分、その文字をファイルに書き出す
				for j := 0; j < count; j++ {
					fmt.Fprintf(file, "%c", currentChr)
				}
				// カッコの分（ {300} など）だけインデックスをスキップ
				i += len(fmt.Sprint(count)) + 2
			} else {
				// 万が一読み取りに失敗したらそのまま1文字として書き出す
				fmt.Fprintf(file, "%c", s[i])
			}
		} else {
			// 普通の文字はそのまま1文字書き出す
			fmt.Fprintf(file, "%c", s[i])
		}
	}
	fmt.Println("展開が完了したなり！")
}
func main() {
  if len(os.Args) < 4 {
    tukaikata()
  }
  switch os.Args[1] {
    case "pack":
      pack(os.Args[2],os.Args[3])
    case "unpack":
      unpack(os.Args[2],os.Args[3])
    default:
      tukaikata()
  }
}
