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
	tokens      []string
}

type Index struct {
	logger    *zap.SugaredLogger
	documents []IndexedDocument
	config    *config.Config
}

func NewIndex(
	config *config.Config,
	logger *zap.SugaredLogger,
) (*Index, error) {
	startTime := time.Now()
	index := &Index{
		config:    config,
		logger:    logger,
		documents: make([]IndexedDocument, 0, config.MaxSearchResults),
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

			index.documents = append(index.documents, index.getTokenizedDocument(record))
		}
	}

	logger.With(
		"duration", time.Since(startTime),
	).Info("Loading index complete...")

	return index, nil
}

func (i *Index) getTokenizedDocument(record []string) IndexedDocument {
	return IndexedDocument{
		Description: record[3],
		ImageUrl:    record[4],
		tokens:      i.getTokens(record[3]),
	}
}

func (i *Index) Search(keywords ...string) ([]IndexedDocument, error) {
	results := make([]IndexedDocument, 0)

	for _, document := range i.documents {
		for _, keyword := range keywords {
			if slices.Index(document.tokens, i.tokenizer(keyword)) >= 0 {
				results = append(results, document)

				break
			}
		}

		if len(results) >= i.config.MaxSearchResults {
			break
		}
	}

	return results, nil
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
