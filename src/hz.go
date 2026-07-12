package main
import (
	"fmt"
	"os"
	"io"
	"strings"
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

	// 【ここから大改造！】
	// カッコ（{）が文字列の中に含まれている限り、何度も繰り返し解凍するループ
	for strings.Contains(s, "{") {
		var nextStr string // 今回の周回で解凍した文字を溜める変数
		
		for i := 0; i < len(s); i++ {
			// もし「文字 + {」を見つけたら、一番内側のカッコを解凍する
			if i+1 < len(s) && s[i+1] == '{' {
				currentChr := s[i]
				var count int
				
				// カッコの中の数字を読み取る
				_, err := fmt.Sscanf(s[i+1:], "{%d}", &count)
				
				if err == nil {
					// 読み取れた数だけ、文字を組み立てる
					for j := 0; j < count; j++ {
						nextStr += string(currentChr)
					}
					// カッコの分（{36}など）だけ進める
					i += len(fmt.Sprint(count)) + 1
				} else {
					nextStr += string(s[i])
				}
			} else {
				nextStr += string(s[i])
			}
		}
		s = nextStr // 解凍途中の文字列を s に上書きして、もう一周チェックする！
	}

	// 最後に、完全に解凍し終わった s をファイルに書き出す
	fmt.Fprint(file, s)
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
