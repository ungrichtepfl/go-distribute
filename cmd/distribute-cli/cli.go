package main

import (
	"fmt"
	"github.com/ungrichtepfl/go-distribute/pkg/excel"
)

func main() {
	last_name_col := uint64(0)
	raw_data, err := excel.EmailToName("./test/data/test_info.xlsx", 1, 3, &last_name_col, nil, nil, nil)
	fmt.Println(raw_data)
	fmt.Println(err)

}
