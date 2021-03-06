package tinysearch

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type IndexWriter struct {
	indexDir string
}

func NewIndexWriter(path string) *IndexWriter {
	return &IndexWriter{
		indexDir: path,
	}
}

// インデックスをファイルに保存する
func (iw *IndexWriter) Flush(index *Index) error {
	for term, postingsList := range index.Dictionary {
		if err := iw.postingsList(term, postingsList); err != nil {
			fmt.Printf("failed to save %s postings list %v", term, err)
		}
	}
	return iw.docCount(index.TotalDocsCount)
}

// ポスティングリストをファイルに保存する
func (iw *IndexWriter) postingsList(term string, postingsList PostingsList) error {
	bytes, err := json.Marshal(postingsList)
	if err != nil {
		return err
	}

	filename := filepath.Join(iw.indexDir, term)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	writer := bufio.NewWriter(file)
	_, err = writer.Write(bytes)
	if err != nil {
		return err
	}

	return writer.Flush()
}

// インデックスされたドキュメント総数をファイルに保存する
func (iw *IndexWriter) docCount(count int) error {
	filename := filepath.Join(iw.indexDir, "_0.dc")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(strconv.Itoa(count)))

	return err
} 