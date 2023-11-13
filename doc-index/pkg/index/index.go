package index

import (
	"encoding/csv"
	"github.com/getnexar/golang-programming-task/doc-index/pkg/config"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"go.uber.org/zap"
)

/*
	Naive implementation of a search index.
	Should be improved.
*/

type IndexedDocument struct {
	Description string `json:"description"`
	ImageUrl    string `json:"imageUrl"`
	deleted     bool
}

type Index struct {
	logger        *zap.SugaredLogger
	documents     []IndexedDocument
	config        *config.Config
	invertedIndex map[string][]int
}

type AffectedDocuments struct {
	Count int `json:"count"`
}

func NewIndex(
	config *config.Config,
	logger *zap.SugaredLogger,
) (*Index, error) {
	startTime := time.Now()
	index := &Index{
		config:        config,
		logger:        logger,
		documents:     make([]IndexedDocument, 0, config.MaxSearchResults),
		invertedIndex: make(map[string][]int),
	}

	dir, err := os.Open(config.IndexDataDir)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	documentIndex := 0
	for _, filename := range files {
		if !strings.HasSuffix(filename, ".csv") {
			// Skip non-csv files
			continue
		}
		filepath := path.Join(config.IndexDataDir, filename)
		logger.
			With("filepath", filepath).
			Info("Loading data file")
		f, err := os.Open(filepath)
		if err != nil {
			logger.
				With("filepath", filepath).
				With("error", err).
				Error("Unable to read input file ")
		}

		csvReader := csv.NewReader(f)

		for {
			record, err := csvReader.Read()
			if err != nil {
				// Make sure we read the whole file
				break
			}

			if len(record) < 5 || record[0] == "AuthorID" {
				// Skip invalid records or header
				continue
			}

			index.documents = append(index.documents, index.getDocument(record))
			index.updateInvertedIndex(index.getTokens(record[3]), documentIndex)

			documentIndex++
		}
	}

	logger.With(
		"duration", time.Since(startTime),
	).Info("Loading index complete...")

	return index, nil
}

func (i *Index) updateInvertedIndex(tokens []string, documentIndex int) {
	for _, token := range tokens {
		if _, exists := i.invertedIndex[token]; !exists {
			i.invertedIndex[token] = []int{documentIndex}
		} else {
			i.invertedIndex[token] = append(i.invertedIndex[token], documentIndex)
		}
	}
}

func (i *Index) getDocument(record []string) IndexedDocument {
	return IndexedDocument{
		Description: record[3],
		ImageUrl:    record[4],
		deleted:     false,
	}
}

func (i *Index) Search(keywords ...string) []IndexedDocument {
	matches := i.query(keywords)

	if len(matches) == 0 {
		return []IndexedDocument{}
	}

	results := make([]IndexedDocument, 0, len(matches))

	for _, documentIndex := range matches {
		results = append(results, i.documents[documentIndex])
	}

	return results
}

func (i *Index) Delete(keywords ...string) AffectedDocuments {
	deletedDocumentsCount := 0

	// i.query(
	// 	keywords, func(documentIndex int, document *IndexedDocument) {
	// 		document.deleted = true
	// 		deletedDocumentsCount++
	// 	},
	// )

	i.purgeDeletedDocuments()

	return AffectedDocuments{
		Count: deletedDocumentsCount,
	}
}

func (i *Index) purgeDeletedDocuments() AffectedDocuments {
	purgedDocumentsCount := 0
	result := make([]IndexedDocument, 0, len(i.documents))

	for _, document := range i.documents {
		if !document.deleted {
			result = append(result, document)
			purgedDocumentsCount++
		}
	}

	i.documents = result

	return AffectedDocuments{
		Count: purgedDocumentsCount,
	}
}

func (i *Index) query(keywords []string) []int {
	matches := i.andMatch(keywords)

	if len(matches) > 0 {
		if len(matches) > i.config.MaxSearchResults {
			return matches[:i.config.MaxSearchResults]
		} else {
			return matches
		}
	}

	return []int{}
}

func (i *Index) andMatch(keywords []string) []int {
	matches := make([][]int, 0, len(keywords))

	for _, keyword := range keywords {
		keyword = i.tokenizer(keyword)

		documentIndexes, exists := i.invertedIndex[keyword]

		if !exists || len(documentIndexes) == 0 {
			return []int{}
		}

		matches = append(matches, documentIndexes)
	}

	return i.intersect(matches[0], matches[1:]...)
}

func (i *Index) intersect(slice1 []int, otherSlices ...[]int) []int {
	result := make([]int, 0, len(slice1))

	for _, element := range slice1 {
		foundInAll := true

		for _, otherSlice := range otherSlices {
			_, found := slices.BinarySearch(otherSlice, element)

			if !found {
				foundInAll = false
				break
			}
		}

		if foundInAll {
			result = append(result, element)
		}
	}

	return result
}

func (i *Index) getTokens(text string) []string {
	tokens := strings.Fields(text)

	result := make([]string, 0, len(tokens))

	for _, token := range tokens {
		token = i.tokenizer(token)

		if token != "" {
			result = append(result, token)
		}
	}

	slices.Sort(result)

	return result
}

func (i *Index) tokenizer(token string) string {
	token = strings.Trim(token, ".,:;!?'\" ")
	token = strings.ToUpper(token)

	if len(token) >= i.config.MinTokenLength && len(token) <= i.config.MaxTokenLength {
		return token
	}

	return ""
}
