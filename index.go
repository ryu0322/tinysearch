package tinysearch

import (
	"container/list"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
	"strconv"
)

// 転置インデックス
type Index struct {
	Dictionary map[string]PostingsList	// 辞書
	TotalDocsCount int					// ドキュメントの総数
}

// New Index create a new Index
func NewIndex() *Index {
	dict := make(map[string]PostingsList)
	return &Index{
		Dictionary: dict,
		TotalDocsCount: 0,
	}
}

// ドキュメントID
type DocumentID int64

// ポスティング
type Posting struct {
	DocID DocumentID	// ドキュメントID
	Positions []int		// 用語の出現位置
	TermFrequency int	// ドキュメント内の用語の出現回数
}

// ポスティングを作成する
func NewPostion(docID DocumentID, positions ...int) *Posting {
	return &Posting {
		docID,
		positions,
		len(positions),
	}
}

// ポスティングリスト
type PostingsList struct {
	*list.List
}

// ポスティングリストを作成する
func NewPostingList (postings ...*Posting) PostingsList {
	li := list.New()
	for _, posting := range postings {
		li.PushBack(posting)
	}
	return PostingsList{li}
}

func (pl PostingsList) add(pos *Posting) {
	pl.PushBack(pos)
}

func (pl PostingsList) last() *Posting {
	e := pl.List.Back()
	if e == nil {
		return nil
	}
	return e.Value.(*Posting)
}

func (idx Index) String() string {
	var padding int
	keys := make([]string, 0, len(idx.Dictionary))
	for key, _ := range idx.Dictionary {
		l := utf8.RuneCountInString(key)
		if padding < l {
			padding = l
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	strs := make([]string, len(keys))
	format := "  [%-" + strconv.Itoa(padding) + "s] -> %s"
	for icnt, key := range keys {
		if postingList, ok := idx.Dictionary[key]; ok {
			strs[icnt] = fmt.Sprintf(format, key, postingList.String())
		}
	}
	return fmt.Sprintf("total documents : %v\ndictionary:\n%v\n",
		idx.TotalDocsCount, strings.Join(strs, "\n"))
}

func (pl *PostingsList) MarshalJson() ([]byte, error) {
	postings := make([]*Posting, 0, pl.Len())

	for e := pl.Front(); e != nil; e = e.Next() {
		postings = append(postings, e.Value.(*Posting))
	}

	return json.Marshal(postings)
}

func (pl *PostingsList) UnMarshalJson(b []byte) error {
	var postings []*Posting
	if err := json.Unmarshal(b, &postings); err != nil {
		return err
	}

	pl.List = list.New()
	for _, posting := range postings {
		pl.add(posting)
	}

	return nil
}

// ポスティングをリストに追加する
// ポスティングリストの最後を取得してドキュメントIDが一致していなければポスティングを追加
// 一致していればポジションを追加
func (pl PostingsList) Add(new *Posting) {
	last := pl.last()
	if last == nil || last.DocID != new.DocID {
		pl.add(new)
		return 
	}
	last.Positions = append(last.Postions, new.Positions)
	last.TermFrequency++
}

func (pl PostingsList) String() string {
	str := make([]string, 0, pl.Len())
	for e := pl.Front(); e != nil; e = e.Next() {
		str = append(str, e.Value.(*Posting).String())
	}
	return strings.Join(str, "=>")
}

// ポスティングリストをたどるためのカーソル
type Cursor struct {
	postingList *PostingsList	// Cursorがたどっているポスティングリストへの参照
	current *list.Element		// 現在の読み込み位置
}

func (pl PostingsList) OpenCursor() *Cursor {
	return &Cursor {
		postingList: &pl,
		current: pl.Front(),
	}
}

func (cur *Cursor) Next() {
	cur.Current = cur.Current.Next()
}

func (cur *Cursor) NextDocId(id DocumentID) {
	for !cur.Empty() && cur.DocId < id {
		cur.Next()
	}
}

func (cur *Cursor) Empty() bool {
	if cur.Current != nil {
		return true
	}

	return false
}

func (cur *Cursor) Posting() *Posting {
	return cur.Current.Value.(*Posting)
}

func (cur *Cursor) DocId() DocumentID {
	return cur.Current.Value.(*Posting).DocID
}

func (cur *Cursor) String() string {
	return fmt.Sprint(cur.Posting())
}