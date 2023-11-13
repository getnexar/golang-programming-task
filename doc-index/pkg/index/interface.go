package index

type IndexInterface interface {
	Search(keywords ...string) []IndexedDocument
	Delete(keywords ...string) AffectedDocuments
}
