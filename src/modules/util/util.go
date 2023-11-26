package util

import (
	"fmt"
	"os"
	types "racc/modules/types"
)

func RemoveDuplicateFilenames(files []types.File) []types.File {
	uniqueNames := make(map[string]struct{})
	result := make([]types.File, 0)

	// Iterate over the files and filter out duplicates
	for _, file := range files {
		// Check if the name is already in the map
		if _, exists := uniqueNames[file.Name]; !exists {
			// Add the file to the result slice and the name to the map
			result = append(result, file)
			uniqueNames[file.Name] = struct{}{}
		}
	}

	return result
}

func RemoveDirectories(files []types.File) []types.File {
	var result []types.File

	for _, file := range files {
		if file.Type != "dir" {
			result = append(result, file)
		}
	}

	return result
}

func EnsureDirExists(dirPath string) error {
	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory: %s", err)
		}
		fmt.Printf("Directory '%s' created\n", dirPath)
	}
	return nil
}
