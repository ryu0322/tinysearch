package tinysearch

import (
	"fmt"
	"math"
	"sort"
)

// searchTopKの検索結果を保持する
type TopDocs struct {
	totalHits int				// ヒット件数
	scoreDocs []*ScoreDoc		// 検索結果
}

func (t *TopDocs) String() string {
	return fmt.Sprintf("\ntotal hits: %v\nresults: %v\n",
			t.totalHits, t.scoreDocs)
}

// ドキュメントIDとそのスコアを保存
type ScoreDoc struct {
	docID DocumentID
	score float64 
}

func (d *ScoreDoc) String() string {
	return fmt.Sprintf("docId: %v, Score: %v", d.docID, d.score)
}

// 検索データに必要なデータを保持する
type Searcher struct {
	indexReader *IndexReader
	cursors []*Cursor
}

func NewSearcher(path string) *Searcher {
	return &Searcher{
		indexReader: NewIndexReader(path),
	}
}

func (sea *Searcher) SearchTopK(query []string, k int) *TopDocs {
	// マッチするドキュメントを抽出、スコアを抽出する
	results := sea.search(query)

	// 結果をスコアの降順でソートする
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	total := len(results)
	if total > k {
		results = results[:k]
	}

	return &TopDocs {
		totalHits: total,
		scoreDocs: results,
	}
}

// 検索を実行し、マッチしたドキュメントをスコア付きで返す
func (sea *Searcher) search(query []string) []*ScoreDoc {
	// クエリに含まれる単語のポスティングリストが
	// 一つも存在しない場合、0件で終了する
	if sea.openCursors(query) == 0 {
		return []*ScoreDoc{}
	}

	// 一番短いポスティングリストを参照するカーソルを選択
	c0 := sea.cursors[0]
	cursors := sea.cursors[1:]

	// 結果を格納する構造体を初期化
	docs := make([]*ScoreDoc, 0)

	for !c0.Empty() {
		var NextDocId DocumentID

		for _, cursor := range cursors {
			// docID以上になるまで読み進める
			if cursor.NextDoc(c0.DocId()); cursor.Empty() {
				return docs
			}
			// docIdが一致しなければ
			if cursor.DocId() != c0.DocId() {
				NextDocId = cursor.DocId()
				break
			}
		}
		if NextDocId > 0 {
			// nextDocId以上になるまで進める
			if c0.NextDoc(NextDocId); c0.Empty() {
				return docs
			}
		} else {
			// 結果を格納する
			docs = append(docs, &ScoreDoc{
				docID: c0.DocId(),
				score: sea.calcScore(),
			})
			c0.Next()
		}
	}
	return docs
}

// 検索に使用するポスティングリストのポインタを取得する
// 作成したカーソルの数を返す
func (sea *Searcher) openCursors(query []string) int {
	// ポスティングリストを取得
	postings := sea.indexReader.postingsLists(query)
	if len(postings) == 0 {
		return 0
	}

	// ポスティングの短い順にソート
	sort.Slice(postings, func(i, j int) bool {
		return postings[i].Len() < postings[j].Len()
	})

	// 各ポスティングリストに対するcursorの取得
	cursors := make([]*Cursor, len(postings))
	for idx, postingList := range postings {
		cursors[idx] = postingList.OpenCursor()
	}
	sea.cursors = cursors
	return len(cursors)
}

// tf-idfスコアを計算する
func (s *Searcher) calcScore() float64 {
	var score float64
	for idx := 0; idx < len(s.cursors); idx++ {
		termFreq := s.cursors[idx].Posting().TermFrequency
		docCount := s.cursors[idx].postingList.Len()
		totalDocCount := s.indexReader.totalDocCount()
		score += calcTF(termFreq) * calIDF(totalDocCount, docCount)
	}
	return score
}

// TFの計算
func calcTF(termCount int) float64 {
	if termCount <= 0 {
		return 0
	}
	return math.Log2(float64(termCount)) + 1
}

// IDFの計算
func calIDF(N, df int) float64 {
	return math.Log2(float64(N) / float64(df))
}