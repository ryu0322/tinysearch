package commands

import (
	"fmt"
	"tinysearch"
	"github.com/urfave/cli"
	"strings"
)

var SearchCommand = cli.Command {
	Name: "search",
	Usage: "search_documents",
	ArgsUsage: "<query>",
	Flags: []cli.Flag {
		cli.IntFlag {
			Name: "number n",
			Value: 10,
			Usage: "",
		},
	},
	Action: search,
}

// 検索を実行するコマンド
func search(con *cli.Context) error {
	if err := checkArgs(con, 1, exactArgs); err != nil {
		return err
	}

	query := con.Args().Get(0)
	result, err := engine.Search(query, con.Int("number"))
	if err != nil {
		return err
	}
	printResult(result)
	return nil
}

// 検索結果を表示する
func printResult(results []*tinysearch.SearchResults) {
	if len(results) == 0 {
		fmt.Println("0 match!!")
		return
	}

	s := make([]string, len(results))
	for idx, result := range results {
		s[idx] = fmt.Sprintf("rank:%3d  score:%4f  title:%s",
			i+1, result.Score, result.Title
		)
	}

	fmt.Println(strings.Join(s, "\n"))
}