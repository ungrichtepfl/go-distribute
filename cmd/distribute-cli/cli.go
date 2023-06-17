package main

import (
	"github.com/ungrichtepfl/go-distribute/pkg/excel"
)

func main() {
	excel.NameToEmail("test/data/test_info.xlsx", 0, 2, nil, nil, nil, nil)
}
