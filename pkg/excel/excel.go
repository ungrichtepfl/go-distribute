package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func NameToEmail(file_path string, first_name_col uint8, email_col uint8, last_name_col *uint8, sheet_name *string, start *uint32, stop *uint32) (*map[string]string, error) {

	rows, err := getRows(file_path, sheet_name, start, stop)

	if err != nil {
		return nil, err
	}

	for i, row := range *rows {

		for j, cell := range row {
			fmt.Print("(", i, j, ") ")
			fmt.Print(cell, "\t")
		}
		fmt.Println()
	}
	return nil, nil
}

func getRows(file_path string, sheet_name *string, start *uint32, stop *uint32) (*[][]string, error) {

	f, err := excelize.OpenFile("test_info.xlsx")
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var excel_sheet_name string
	if sheet_name == nil {
		excel_sheet_name = f.GetSheetName(0)
	} else {
		excel_sheet_name = *sheet_name
	}

	if excel_sheet_name == "" {
		return nil, fmt.Errorf("No sheets in excel file found")
	}
	rows, err := f.GetRows(excel_sheet_name)

	if err != nil {
		return nil, err
	}

	if stop != nil {
		rows = rows[:*stop]
	}

	if start != nil {
		rows = rows[*start:]
	}

	return &rows, nil

}
