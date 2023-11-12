package index

type IndexInterface interface {
	Search(keywords ...string) ([]IndexedDocument, error)
	Delete(keywords ...string) AffectedDocuments
}
