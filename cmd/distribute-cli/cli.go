package main

import (
	"fmt"
	"github.com/ungrichtepfl/go-distribute/pkg/distribute"
)

func main() {
	last_name_col := uint64(0)
	file_path := `test/data/test_info.xlsx`
	println(file_path)
	raw_data, err := distribute.EmailToName(file_path, 1, 3, &last_name_col, nil, nil, nil)
	fmt.Println(raw_data)
	fmt.Println(err)

	dir := "test/data/documents"
	files, err := distribute.ListPDFs(dir, false)
	fmt.Println(files)
	fmt.Println(err)
	files, err = distribute.ListPDFs(dir, true)
	fmt.Println(files)
	fmt.Println(err)

}
