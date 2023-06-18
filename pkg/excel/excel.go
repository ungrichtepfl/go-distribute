package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

func EmailToName(file_path string, first_name_col uint64, email_col uint64, last_name_col *uint64, sheet_name *string, start *uint64, stop *uint64) (*map[string]string, *[]string, error) {

	rows, stride_from_document_start, err := getRows(file_path, sheet_name, start, stop)
	if err != nil {
		return nil, nil, err
	}

	return getEmailNameMap(rows, first_name_col, email_col, last_name_col, stride_from_document_start)

}

func getEmailNameMap(rows [][]string, first_name_col uint64, email_col uint64, last_name_col *uint64, stride_from_start uint64) (*map[string]string, *[]string, error) {

	if int(first_name_col) > len(rows[0]) {
		return nil, nil, fmt.Errorf("First name column index out of range. got %d, max %d. Row: %v", first_name_col, len(rows[0]), rows[0])
	}
	if int(email_col) > len(rows[0]) {
		return nil, nil, fmt.Errorf("Email column index out of range. got %d, max %d", email_col, len(rows[0]))
	}
	if last_name_col != nil && int(*last_name_col) > len(rows[0]) {
		return nil, nil, fmt.Errorf("Last name column index out of range. got %d, max %d", *last_name_col, len(rows[0]))
	}

	email_name_map := make(map[string]string, len(rows))

	var names_with_no_emails []string

	for i, row := range rows {
		first_name := strings.TrimSpace(row[first_name_col])
		email := strings.TrimSpace(row[email_col])
		last_name := ""
		if last_name_col != nil {
			last_name = strings.TrimSpace(row[*last_name_col])
		}
		var name string
		if last_name != "" {
			name = first_name + " " + last_name
		} else {
			name = first_name
		}
		if name == "" && email == "" {
			fmt.Printf("No name and email for some row %d found found. Skipping row! Row: %v.\n", stride_from_start+uint64(i), row)
			continue
		}
		if email == "" {
			names_with_no_emails = append(names_with_no_emails, name)
			continue
		}
		email_name_map[name] = email

	}
	return &email_name_map, &names_with_no_emails, nil
}

func getRows(file_path string, sheet_name *string, start *uint64, stop *uint64) ([][]string, uint64, error) {

	f, err := excelize.OpenFile(file_path)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	excel_sheet_name, err := getSheetName(f, sheet_name)
	if err != nil {
		return nil, 0, err
	}

	rows, err := f.GetRows(excel_sheet_name)
	if err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, 0, fmt.Errorf("No data found in sheet %s", excel_sheet_name)
	}

	if stop != nil {
		rows = rows[:*stop]
	}

	table_width := maxRowLength(rows)

	var stride_from_document_start uint64
	if start != nil {
		stride_from_document_start = *start
	} else {
		var err error
		stride_from_document_start, err = findStart(rows, table_width)
		if err != nil {
			return nil, 0, err
		}
	}
	rows = rows[stride_from_document_start:]

	makeRowsConsistent(rows, table_width)

	printRows(rows)

	return rows, stride_from_document_start, nil

}

func makeRowsConsistent(rows [][]string, length uint64) error {

	for i, row := range rows {
		if uint64(len(row)) < length {
			rows[i] = append(row, make([]string, length-uint64(len(row)))...)
		} else if uint64(len(row)) > length {
			panic(fmt.Sprintf("Row %d has length %d, which is greater than the length specified in the arguments (%d)", i, len(row), length))
		}
	}
	return nil

}

func printRows(rows [][]string) {
	fmt.Println("Printing rows:")
	for i, row := range rows {
		for j, col := range row {
			fmt.Printf("(%d, %d) %s\t", i, j, col)
		}
		fmt.Println()
	}

}

func maxRowLength(rows [][]string) uint64 {
	var max_row_length uint64
	for _, row := range rows {
		if uint64(len(row)) > max_row_length {
			max_row_length = uint64(len(row))
		}
	}
	return max_row_length

}

func findStart(rows [][]string, table_width uint64) (uint64, error) {

	for i, row := range rows {
		if uint64(len(row)) == table_width {
			return uint64(i + 1), nil // Do not return i because we want to skip the table header
		}
	}

	return 0, fmt.Errorf("No start found. Please specify start row.")

}

func findStop(rows [][]string, table_width uint64) (uint64, error) {

	for i := range rows {
		row := rows[len(rows)-1-i]
		if uint64(len(row)) == table_width {
			return uint64(len(rows) - i), nil
		}
	}

	return 0, fmt.Errorf("No stop found. Please specify stop row.")

}

func getSheetName(f *excelize.File, sheet_name *string) (string, error) {

	var excel_sheet_name string
	if sheet_name == nil || *sheet_name == "" {
		excel_sheet_name = f.GetSheetName(0)
	} else {
		excel_sheet_name = *sheet_name
	}

	if excel_sheet_name == "" {
		return "", fmt.Errorf("No sheets in excel file found")
	}

	return excel_sheet_name, nil
}
