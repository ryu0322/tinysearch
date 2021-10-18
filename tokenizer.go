package tinysearch

import (
	"bufio"
	"bytes"
	"strings"
	"unicode"
)

type Tokenizer struct {}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

func replace (rep rune) rune {
	// 英数字以外は除外
	if (rep < 'a' || rep > 'z') && (rep < 'A' || rep > 'Z') && !unicode.IsNumber(rep) {
		return -1
	}

	// 大文字を小文字に変換
	return unicode.ToLower(rep)
}

// io.Readerから読んだデータをトークンに分割する
func (t *Tokenizer) SplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanWords(data, atEOF)
	if err == nil && token != nil {
		token = bytes.Map(replace, token)
		if len(token) == 0 {
			token = nil
		}
	}

	return
}

// 文字列を分解
func (t *Tokenizer) TextToWordSequence(text string) []string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(t.SplitFunc)
	var result []string
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result
}