package index

import (
	"github.com/getnexar/golang-programming-task/doc-index/pkg/config"
	"go.uber.org/zap"
	"os"
	"path"
	"reflect"
	"slices"
	"testing"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		expectedToken string
	}{
		{"Valid Token", "Hello", "HELLO"},
		{"Empty Token", "", ""},
		{"Too Short Token", "ab", ""},
		{"Too Long Token", "toolongtoken", ""},
		{"Valid Token ends with !", "word!", "WORD"},
		{"Valid Token ends with .", "word.", "WORD"},
		{"Valid Token ends with ,", "word,", "WORD"},
		{"Valid Token ends with :", "word:", "WORD"},
		{"Valid Token ends with ;", "word;", "WORD"},
		{"Valid Token ends with ?", "word?", "WORD"},
		{"Valid Token ends with '", "word'", "WORD"},
		{"Valid Token ends with \"", "word\"", "WORD"},
		{"Valid Token starts with !", "!word", "WORD"},
		{"Valid Token starts with .", ".word", "WORD"},
		{"Valid Token starts with ,", ",word", "WORD"},
		{"Valid Token starts with :", ":word", "WORD"},
		{"Valid Token starts with ;", ";word", "WORD"},
		{"Valid Token starts with ?", "?word", "WORD"},
		{"Valid Token starts with '", "'word", "WORD"},
		{"Valid Token starts with \"", "\"word", "WORD"},
		{"Valid Token with : inside", "wor:d", "WOR:D"},
		{"Valid Token with Whitespace", "   Spaces   ", "SPACES"},
	}

	index := newIndex()

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				result := index.tokenizer(test.token)
				if result != test.expectedToken {
					t.Errorf("Expected token %s, got %s", test.expectedToken, result)
				}
			},
		)
	}
}

func TestGetTokens(t *testing.T) {
	tests := []struct {
		name         string
		inputText    string
		expectedList []string
	}{
		{"Simple Case", "Hello World", []string{"HELLO", "WORLD"}},
		{"Empty String", "", []string{}},
		{"String with Punctuation", "word! token.", []string{"WORD", "TOKEN"}},
		{"String with Whitespace", "  Spaces  Between  ", []string{"SPACES", "BETWEEN"}},
		{"Too long token", "Hello Toolongtoken", []string{"HELLO"}},
	}

	index := newIndex()

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				result := index.getTokens(test.inputText)

				slices.Sort(test.expectedList)

				if !reflect.DeepEqual(result, test.expectedList) {
					t.Errorf("Expected tokens %v, got %v", test.expectedList, result)
				}
			},
		)
	}
}

func TestGetTokenizedDocument(t *testing.T) {
	index := newIndex()

	tests := []struct {
		name             string
		inputRecord      []string
		expectedDocument IndexedDocument
	}{
		{
			"Simple Case",
			[]string{"field1", "field2", "field3", "Simple Document", "image-url"},
			IndexedDocument{
				Description: "Simple Document",
				ImageUrl:    "image-url",
				tokens:      index.getTokens("Simple Document"),
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				result := index.getTokenizedDocument(test.inputRecord)

				if !reflect.DeepEqual(result, test.expectedDocument) {
					t.Errorf("Expected document %v, got %v", test.expectedDocument, result)
				}
			},
		)
	}
}

func TestNewIndex(t *testing.T) {
	configOptions, _ := config.Load()

	logger := zap.NewNop().Sugar()

	t.Run(
		"ValidData", func(t *testing.T) {
			tempDir := os.TempDir()
			testCSVPath := path.Join(tempDir, "test_data_1.csv")
			testSecondCSVPath := path.Join(tempDir, "test_data_2.csv")
			testNonCSVPath := path.Join(tempDir, "test_data.txt")

			defer os.Remove(testCSVPath)
			defer os.Remove(testNonCSVPath)
			defer os.Remove(testSecondCSVPath)

			createTestFile(
				t, testCSVPath, []byte(
					"AuthorID,Author,Date,Content,Attachments\n"+
						"\"1\",\"Author\",\"04/20/2023 12:00 AM\",\"content\",\"image-url\"\n"+
						"\"2\",\"Author\",\"04/20/2023 12:00 AM\",\"invalid\"\n",
				),
			)

			createTestFile(
				t, testSecondCSVPath, []byte(
					"AuthorID,Author,Date,Content,Attachments\n"+
						"\"1\",\"Author\",\"04/20/2023 12:00 AM\",\"valid\",\"image-url\"\n",
				),
			)

			createTestFile(
				t, testNonCSVPath, []byte(
					"AuthorID,Author,Date,Content,Attachments\n"+
						"\"1\",\"Author\",\"04/20/2023 12:00 AM\",\"content\",\"image-url\"\n"+
						"\"2\",\"Author\",\"04/20/2023 12:00 AM\",\"invalid\"\n",
				),
			)

			configOptions.IndexDataDir = tempDir

			index, err := NewIndex(configOptions, logger)

			if err != nil {
				t.Errorf("Error creating index: %v", err)
			}

			if len(index.documents) != 2 {
				t.Errorf("Expected index to contain 2 documents, got %d", len(index.documents))
			}

			if index.documents[0].Description != "content" {
				t.Errorf("Expected document description to be 'content', got %s", index.documents[0].Description)
			}

			if index.documents[0].ImageUrl != "image-url" {
				t.Errorf("Expected document image url to be 'image-url', got %s", index.documents[0].ImageUrl)
			}
		},
	)
}

func TestQuery(t *testing.T) {
	index := newIndex()
	index.documents = getSampleData(index)
	index.config.MaxSearchResults = len(index.documents)

	tests := []struct {
		name              string
		keywords          []string
		expectedDocuments []IndexedDocument
	}{
		{
			name:              "First Document",
			keywords:          []string{"one"},
			expectedDocuments: []IndexedDocument{index.documents[0]},
		},
		{
			name:              "Second document. Case insensitive. Two keywords",
			keywords:          []string{"two", "Document"},
			expectedDocuments: []IndexedDocument{index.documents[1]},
		},
		{
			name:              "Query for three keywords, None match",
			keywords:          []string{"one", "document", "four"},
			expectedDocuments: []IndexedDocument{},
		},
		{
			name:              "All documents match",
			keywords:          []string{"document"},
			expectedDocuments: index.documents,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				results := make([]IndexedDocument, 0)

				index.query(
					test.keywords, func(documentIndex int, document *IndexedDocument) {
						results = append(results, *document)
					},
				)

				if !reflect.DeepEqual(results, test.expectedDocuments) {
					t.Errorf("Expected documents %v, got %v", test.expectedDocuments, results)
				}
			},
		)
	}
}

func TestMaxSearchResults(t *testing.T) {
	index := newIndex()
	index.documents = getSampleData(index)

	tests := []struct {
		name       string
		keywords   []string
		maxResults int
	}{
		{
			name:       "3 documents match",
			keywords:   []string{"document"},
			maxResults: 3,
		},
		{
			name:       "2 documents match",
			keywords:   []string{"document"},
			maxResults: 2,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				foundDocuments := 0

				index.config.MaxSearchResults = test.maxResults

				index.query(
					test.keywords, func(documentIndex int, document *IndexedDocument) {
						foundDocuments++
					},
				)

				if foundDocuments != test.maxResults {
					t.Errorf("Expected %d results, got %d", test.maxResults, foundDocuments)
				}
			},
		)
	}
}

func TestSearch(t *testing.T) {
	index := newIndex()
	index.documents = getSampleData(index)
	index.config.MaxSearchResults = 10

	tests := []struct {
		name              string
		keywords          []string
		expectedDocuments []IndexedDocument
	}{
		{
			name:              "First Document",
			keywords:          []string{"one"},
			expectedDocuments: []IndexedDocument{index.documents[0]},
		},
		{
			name:              "Second document. Case insensitive. Two keywords",
			keywords:          []string{"two", "Document"},
			expectedDocuments: []IndexedDocument{index.documents[1]},
		},
		{
			name:              "Query for three keywords, None match",
			keywords:          []string{"one", "document", "four"},
			expectedDocuments: []IndexedDocument{},
		},
		{
			name:              "All documents match",
			keywords:          []string{"document"},
			expectedDocuments: index.documents,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				results, err := index.Search(test.keywords...)

				if err != nil {
					t.Errorf("Error searching index: %v", err)
				}

				if !reflect.DeepEqual(results, test.expectedDocuments) {
					t.Errorf("Expected documents %v, got %v", test.expectedDocuments, results)
				}
			},
		)
	}
}

func TestDelete(t *testing.T) {
	index := newIndex()
	index.documents = getSampleData(index)
	index.config.MaxSearchResults = 10

	tests := []struct {
		name                   string
		keywords               []string
		affectedDocumentsCount int
		expectedDocuments      []IndexedDocument
	}{
		{
			name:                   "Delete First Document",
			keywords:               []string{"one"},
			affectedDocumentsCount: 1,
			expectedDocuments:      []IndexedDocument{index.documents[1], index.documents[2]},
		},
		{
			name:                   "Delete First and Second documents",
			keywords:               []string{"top", "secret"},
			affectedDocumentsCount: 2,
			expectedDocuments:      []IndexedDocument{index.documents[2]},
		},
		{
			name:                   "No documents to delete",
			keywords:               []string{"document", "not", "found"},
			affectedDocumentsCount: 0,
			expectedDocuments:      index.documents,
		},
	}

	for _, test := range tests {
		// Restore deleted documents
		index.documents = getSampleData(index)

		t.Run(
			test.name, func(t *testing.T) {
				affectedDocuments := index.Delete(test.keywords...)

				if affectedDocuments.Count != test.affectedDocumentsCount {
					t.Errorf(
						"Expected %d affected documents, got %d", test.affectedDocumentsCount, affectedDocuments.Count,
					)
				}

				if !reflect.DeepEqual(index.documents, test.expectedDocuments) {
					t.Errorf("Expected documents %v, got %v", test.expectedDocuments, index.documents)
				}
			},
		)
	}
}

func newIndex() *Index {
	configOptions, _ := config.Load()

	configOptions.MaxTokenLength = 10

	return &Index{
		config: configOptions,
	}
}

func createTestFile(t *testing.T, filename string, data []byte) {
	err := os.WriteFile(filename, data, 0644)

	if err != nil {
		t.Errorf("Error writing non CVS file: %v", err)
	}
}

func getSampleData(index *Index) []IndexedDocument {
	return []IndexedDocument{
		{
			Description: "Document One. Top secret!",
			ImageUrl:    "image-url",
			tokens:      index.getTokens("Document One. Top secret!"),
		},
		{
			Description: "Document Two. Top secret!",
			ImageUrl:    "image-url",
			tokens:      index.getTokens("Document Two. Top secret!"),
		},
		{
			Description: "Document Three",
			ImageUrl:    "image-url",
			tokens:      index.getTokens("Document Three"),
		},
	}
}
