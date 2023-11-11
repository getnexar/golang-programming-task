package index

import (
	"encoding/csv"
	"os"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
)

/*
	Naive implementation of a search index.
	Should be improved.
*/

type Index struct {
	logger     *zap.SugaredLogger
	documents  [][2]string
	maxResults int
}

func NewIndex(
	indexDataDir string,
	maxResults int,
	logger *zap.SugaredLogger,
) (*Index, error) {
	startTime := time.Now()
	index := &Index{
		maxResults: maxResults,
		logger:     logger,
		documents:  make([][2]string, 0),
	}
	dir, err := os.Open(indexDataDir)
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
		filepath := path.Join(indexDataDir, filename)
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
			if len(record) < 5 {
				// Skip invalid records
				continue
			}
			index.documents = append(index.documents, [2]string{record[3], record[4]})
		}
	}

	logger.With(
		"duration", time.Since(startTime),
	).Info("Loading index complete...")

	return index, nil
}

func (i *Index) Search(keywords ...string) ([][2]string, error) {
	results := make([][2]string, 0)
	for _, document := range i.documents {
		found := true
		for _, keyword := range keywords {
			if !strings.Contains(document[0], keyword) {
				found = false
				break
			}
		}
		if found {
			results = append(results, document)
		}
		if len(results) >= i.maxResults {
			break
		}
	}
	return results, nil
}
