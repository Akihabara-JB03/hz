package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func tukaikata() {
	fmt.Println("使い方をしれよ！！")
	fmt.Println("使い方は、hz [pack/unpack] ファイル名.その拡張子 変換するファイル名.hz(展開のときは、ファイル名.hz ファイル名.さっきまでの拡張子)")
	os.Exit(1)
}

// 【圧縮機能】同じ文字が5回以上続いたら A{5} の形に変換するなり！
func pack(input string, output string) {
	file, err := os.Create(output)
	if err != nil {
		fmt.Printf("圧縮ファイル作成エラーなり！！　詳細はこれなり！！ エラー内容:%s\n", err)
		return
	}
	defer file.Close()

	inpfile, err := os.Open(input)
	if err != nil {
		fmt.Printf("圧縮元ファイル確認準備エラーなり！！　詳細はこれなり！！　エラー内容:%s\n", err)
		return
	}
	defer inpfile.Close()

	data, err := io.ReadAll(inpfile)
	if err != nil {
		fmt.Printf("圧縮元ファイル確認エラーなり！！詳細はこれなり！！　エラー内容:%s\n", err)
		return
	}

	// 日本語（マルチバイト文字）が混ざってもバイト数がズレないように rune 配列に変換するなり！
	s := []rune(string(data))
	if len(s) == 0 {
		fmt.Println("空ファイルなり！圧縮完了なり！！")
		return
	}

	count := 1
	for i := 0; i < len(s); i++ {
		if i+1 < len(s) && s[i] == s[i+1] {
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

// 【展開機能】A{5} のようなフォーマットを元の文字に完全に復元するなり！
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

	s := []rune(string(data))
	var result strings.Builder // 結合を爆速にするバッファ

	for i := 0; i < len(s); i++ {
		// もし「{」を見つけたら、直前の文字が「繰り返す文字」なり！
		if s[i] == '{' {
			// 直前の文字を取得（すでにresultに書き込まれているので、1文字削って使い回す）
			strSoFar := result.String()
			runesSoFar := []rune(strSoFar)
			if len(runesSoFar) == 0 {
				// 直前に文字がない不正なデータはそのまま出力してスキップ
				result.WriteRune(s[i])
				continue
			}
			targetChar := runesSoFar[len(runesSoFar)-1]

			// resultの末尾から、今コピー対象にした1文字を削除する
			result.Reset()
			result.WriteString(string(runesSoFar[:len(runesSoFar)-1]))

			// カッコの中の数字（リピート回数）を読み取る
			count := 0
			numLen := 0
			hasClosing := false
			for j := i + 1; j < len(s); j++ {
				if s[j] >= '0' && s[j] <= '9' {
					count = count*10 + int(s[j]-'0')
					numLen++
				} else if s[j] == '}' {
					hasClosing = true
					break
				} else {
					break
				}
			}

			// 正常に「}」で閉じられていたら展開するなり！
			if hasClosing && numLen > 0 {
				for k := 0; k < count; k++ {
					result.WriteRune(targetChar)
				}
				i += numLen + 1 // 数字の長さ ＋ 「}」の分、インデックスを進める
			} else {
				// カッコのパースに失敗した場合は、削った文字を戻して普通に処理する
				result.WriteRune(targetChar)
				result.WriteRune(s[i])
			}
			continue
		}

		// 普通の文字はそのままバッファに突っ込むなり！
		result.WriteRune(s[i])
	}

	fmt.Fprint(file, result.String())
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
