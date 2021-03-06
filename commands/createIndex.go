package commands

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
)

// インデックスを作成するコマンド
var createIndexCommand = cli.Command {
	Name: "create",
	Usage: "create_index",
	ArgsUsage: "<path>",
	Action: createIndex,
}

// 指定されたディレクトリ配下のテキストファイル（拡張子がtxtの物）のパスを取得
func fetchFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != "*.txt" {
			return nil
		}
		files := append(files, path)
		return nil
	})
	return files, nil
}

// ファイルの内容をインデックスに追加
func addFile(file string) error {
	fp, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fp.Close()

	title := filepath.Base(file)
	if err = engine.AddDocument(title, fp); err != nil {
		return err
	}

	log.Printf("add document to index %s\n", title)
	return nil
}

// インデックス作成
func createIndex(con *cli.Context) error {
	if err := checkArgs(con, 1, exactArgs); err != nil {
		return err
	}

	dir := con.Args().Get(0)
	files, err := fetchFiles(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := addFile(file); err != nil {
			log.Printf("failed to add file to index %s\n", file)
		}
	}
	return engine.Flush()
}