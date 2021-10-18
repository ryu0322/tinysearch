package tinysearch

import (
	"database/sql"
	"log"
)

type DocumentStore struct {
	db *sql.DB
}

func NewDocumentStore(db *sql.DB) *DocumentStore {
	return &DocumentStore {
		db: db,
	}
}

func (ds *DocumentStore) save(title string) (DocumentID, error) {
	query := "insert into documents (document_title) values(?)"
	result, err := ds.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	return DocumentID(id), nil

}

func (ds *DocumentStore) fetchTitle(docID DocumentID) (string, error) {
	query := "select document_title from documents where document_id = ?"
	row := ds.db.QueryRow(query, docID)
	
	var title string
	err := row.Scan(&title)
	return title, err
}