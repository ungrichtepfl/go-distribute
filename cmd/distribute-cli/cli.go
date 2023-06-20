package main

import (
	"fmt"
	"github.com/ungrichtepfl/go-distribute/pkg/distribute"
)

func main() {
	last_name_col := uint64(0)
	file_path := "test/data/test_info.xlsx"
	document_dir := "test/data/documents"
	excel_config := distribute.ExcelConfig{
		FilePath:     file_path,
		FirstNameCol: 1,
		LastNameCol:  &last_name_col,
		EmailCol:     3,
		SheetName:    nil,
		StartRow:     nil,
		EndRow:       nil,
	}
	config := distribute.Config{
		ExcelConfig:     excel_config,
		DocumentDir:     document_dir,
		FileTypes:       []distribute.FileType{distribute.PDF, distribute.TXT},
		RecursiveSearch: true,
	}
	lexed_data, err := distribute.EmailToName(config.ExcelConfig)

	fmt.Println("Lexed Data:\n", lexed_data)
	fmt.Println("Error:\n", err)

	fmt.Println("------------------------------------------------------")

	parsed_data, err := distribute.NameEmailToDocuments(config, lexed_data.NameEmails)
	fmt.Println("Parsed Data\n", parsed_data)
	fmt.Println("Error:\n", err)

}
