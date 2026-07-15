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
// 展開用の一時データを保持するスタックの構造体
type unpackTask struct {
	char  rune            // 繰り返したい文字
	count int             // 繰り返す回数
	buf   strings.Builder // カッコの中に溜まった文字列
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

	// Goの文字（UTF-8）を安全にパースするため rune 配列に変換
	s := []rune(string(data))

	// 最外層の文字列をビルドするメインバッファ
	var rootBuf strings.Builder
	// カッコの階層を管理するスタック
	var stack []*unpackTask

	for i := 0; i < len(s); i++ {
		// カッコの開始「{」を見つけた場合
		if s[i] == '{' {
			// 直前の文字が「繰り返す対象の文字」になる
			// スタックが空なら rootBuf から、スタックがあるなら現在のタスクのバッファから1文字剥ぎ取る
			var targetChar rune
			if len(stack) == 0 {
				strSoFar := rootBuf.String()
				if len(strSoFar) > 0 {
					runesSoFar := []rune(strSoFar)
					targetChar = runesSoFar[len(runesSoFar)-1]
					// 剥ぎ取った分、削る
					rootBuf.Reset()
					rootBuf.WriteString(string(runesSoFar[:len(runesSoFar)-1]))
				}
			} else {
				currentTask := stack[len(stack)-1]
				strSoFar := currentTask.buf.String()
				if len(strSoFar) > 0 {
					runesSoFar := []rune(strSoFar)
					targetChar = runesSoFar[len(runesSoFar)-1]
					currentTask.buf.Reset()
					currentTask.buf.WriteString(string(runesSoFar[:len(runesSoFar)-1]))
				}
			}

			// カッコの中の数字（回数）をパースする
			var count int
			numLen := 0
			// 「{」の次の文字から数字が続く限り読み進める
			for j := i + 1; j < len(s); j++ {
				if s[j] >= '0' && s[j] <= '9' {
					count = count*10 + int(s[j]-'0')
					numLen++
				} else {
					break
				}
			}

			// 新しい階層（タスク）を作成してスタックに積む（Push）
			newTask := &unpackTask{
				char:  targetChar,
				count: count,
			}
			stack = append(stack, newTask)

			// 読み進めた数字の分だけインデックスを進める
			i += numLen
			continue
		}

		// カッコの閉じ「}」を見つけた場合
		if s[i] == '}' {
			if len(stack) > 0 {
				// 現在の階層をポップ（Pop）
				finishedTask := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				// カッコ内の文字列（中身）を展開する
				var expanded strings.Builder
				content := finishedTask.buf.String()

				// 「指定文字」を指定回数分ループし、その中にカッコ内の中身を挟み込む
				for k := 0; k < finishedTask.count; k++ {
					expanded.WriteRune(finishedTask.char)
					expanded.WriteString(content)
				}

				// 展開した結果を、1つ上の階層のバッファ（無ければrootBuf）に流し込む
				if len(stack) == 0 {
					rootBuf.WriteString(expanded.String())
				} else {
					stack[len(stack)-1].buf.WriteString(expanded.String())
				}
			}
			continue
		}

		// 通常の文字は、現在の階層のバッファ（無ければrootBuf）にそのまま書き込む
		if len(stack) == 0 {
			rootBuf.WriteRune(s[i])
		} else {
			stack[len(stack)-1].buf.WriteRune(s[i])
		}
	}

	// 最後に、完全に解凍し終わった文字列をファイルに書き出す
	fmt.Fprint(file, rootBuf.String())
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
