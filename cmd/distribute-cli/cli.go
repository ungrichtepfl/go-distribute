package main

import (
	"fmt"
	"github.com/ungrichtepfl/go-distribute/pkg/excel"
)

func main() {
	last_name_col := uint64(0)
	name_to_email_map, names_with_no_emails, err := excel.EmailToName("./test/data/test_info.xlsx", 1, 3, &last_name_col, nil, nil, nil)
	fmt.Println(name_to_email_map)
	fmt.Println(names_with_no_emails)
	fmt.Println(err)

}
