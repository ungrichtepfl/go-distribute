package distribute

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ListByFileType(path string, file_types []FileType, recursive bool) ([]string, error) {

	all_files, err := listDirectory(path, recursive)
	if err != nil {
		return nil, err
	}

	filtered_files := make([]string, 0, len(all_files))
	for _, file := range all_files {
		ext := strings.ToLower(filepath.Ext(file))
		for _, extension := range file_types {
			if ext == extension {
				filtered_files = append(filtered_files, file)
				break
			}
		}
	}
	return filtered_files, nil
}

func listDirectory(path string, recursive bool) ([]string, error) {

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path %s is not a directory", path)
	}

	sep_count := strings.Count(path, string(filepath.Separator))
	depth := func(path string) int {
		return strings.Count(path, string(filepath.Separator)) - sep_count
	}

	var files []string
	const max_depth = 0
	err = filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if !recursive && depth(path) > max_depth {
				return filepath.SkipDir
			}
			return nil
		}
		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in %s", path)
	}

	return files, nil

}
