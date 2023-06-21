package main

import (
	"fmt"
	"github.com/ungrichtepfl/go-distribute/pkg/distribute"
	"os"
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
	email_config := distribute.EmailConfig{
		SenderEmail: "christoph.ungricht@outlook.com",
		SMTPHost:    "smtp.office365.com",
		SMTPPort:    "587",
		Username:    "christoph.ungricht@outlook.com",
		Password:    os.Getenv("EMAIL_PASSWORD"), // App Password
	}

	config := distribute.Config{
		ExcelConfig:     excel_config,
		EmailConfig:     email_config,
		DocumentDir:     document_dir,
		FileTypes:       []distribute.FileType{distribute.PDF, distribute.TXT},
		RecursiveSearch: true,
	}

	err := distribute.SendEmail(config, "christoph.ungricht@outlook.com", []string{})
	if err != nil {
		fmt.Println("Error:\n", err)
		return
	}

	lexed_data, err := distribute.EmailToName(config.ExcelConfig)

	fmt.Println("Lexed Data:\n", lexed_data)
	fmt.Println("Error:\n", err)

	fmt.Println("------------------------------------------------------")

	parsed_data, err := distribute.NameEmailToDocuments(config, lexed_data.NameEmails)
	fmt.Println("Parsed Data\n", parsed_data)
	fmt.Println("Error:\n", err)

}
