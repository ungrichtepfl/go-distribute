package distribute

import (
	"log"
	"path/filepath"
	"strings"
)

func NameEmailToDocuments(config Config, name_email []NameEmail) (ParsedData, error) {
	documents, err := ListByFileType(config.DocumentDir, config.FileTypes, config.RecursiveSearch)
	if err != nil {
		return ParsedData{}, err
	}

	email_to_documents := make(map[NameEmail][]DocumentPath, len(name_email))
	processed_documents := make([]DocumentPath, 0, len(documents))
	anomalies := make([]ParserAnomaly, 0, 5)
	for _, name_email := range name_email {
		documents_found := make([]string, 0, 1)
		for _, document := range documents {
			if isNameInDocumentPath(name_email.Name, document) {
				documents_found = append(documents_found, document)
				processed_documents = append(processed_documents, document)
			}
		}
		if len(documents_found) == 0 {
			log.Printf("No document found for name %s.\n", name_email.Name)
			anomalies = append(anomalies, ParserAnomaly{NoDocumentMatchForName, []string{name_email.Name}})
		} else {
			email_to_documents[name_email] = documents_found
		}
	}

	unprocessed_documents := make([]DocumentPath, 0, len(documents)-len(processed_documents))
	for _, document := range documents {
		found := false
		for _, processed_document := range processed_documents {
			if strings.EqualFold(document, processed_document) {
				found = true
				break
			}
		}
		if !found {
			unprocessed_documents = append(unprocessed_documents, document)
		}
	}

	if len(unprocessed_documents) > 0 {
		log.Printf("No name match found for some documents %v.\n", unprocessed_documents)
		anomalies = append(anomalies, ParserAnomaly{NoNameMatchForDocument, unprocessed_documents})
	}

	return ParsedData{email_to_documents, anomalies}, nil

}

func isNameInDocumentPath(name string, path string) bool {
	name = strings.ToLower(name)
	document_name := strings.ToLower(filepath.Base(path))

	minimum_matches := 2
	total_matches := 0

	for _, word := range strings.Split(name, " ") {
		if strings.Contains(document_name, word) {
			total_matches++
		}
		if total_matches >= minimum_matches {
			return true
		}

	}
	return false

}
