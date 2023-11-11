package index

type IndexInterface interface {
	Search(keywords ...string) ([][2]string, error)
}
