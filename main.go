package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
	"io/ioutil"
	"bufio"
	"bytes"
	"sort"
)

type Snapshot struct {
	ID string `json:"id"`
}

type SnapshotInfo struct {
	Path  string `json:"path"`
	MTime string `json:"mtime"`
}

type SnapshotEntry struct {
	Path string `json:"path"`
}

func main() {
	//#region CLI param definitions
	s3URL := flag.String("r", "", "Target repository URI")
	secretFilePath := flag.String("sf", "", "File containing repository secret")
	secretDirectInput := flag.String("s", "", "Repository secret")
	flag.Parse()
	//#endregion

	//#region CLI param validation
	if *s3URL == "" {
		fmt.Println("Error: Repository URI is required.")
		return
	}

	if *secretFilePath == "" && *secretDirectInput == "" {
		fmt.Println("Error: Repository secret file OR repository secret flag is required.")
		return
	}

	if *secretFilePath != "" && *secretDirectInput != "" {
		fmt.Println("Error: Repository secret file and repository secret flag can not be used at the same time.")
		return
	}
	//#endregion

	//#region Get repo secret
	var repoSecret string
	if *secretFilePath != "" {
		repoSecretContent, repoSecretErr := ioutil.ReadFile(*secretFilePath)
		if repoSecretErr != nil {
			fmt.Println(repoSecretErr)
		}
	
		repoSecret = string(repoSecretContent)
	} else {
		repoSecret = *secretDirectInput
	}
	//#endregion

	//#region Fetch snapshot list
	fmt.Println(repoSecret);
	cmd := exec.Command("restic", "snapshots", "-r", *s3URL, "-p"+*secretFilePath, "--json")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error fetching snapshot list:", err)
		return
	}

	var snapshots []Snapshot
	err = json.Unmarshal(output, &snapshots)
	if err != nil {
		fmt.Println("Error parsing snapshot list:", err)
		return
	}
	fmt.Println("Found " + fmt.Sprintf("%d", len(snapshots)) + " snapshots.")
	//#endregion

	var accumulatedPaths []string

	for _, snapshot := range snapshots {
		cmd = exec.Command("restic", "ls", snapshot.ID, "-r", *s3URL, "-p"+*secretFilePath, "--json")
		output, err = cmd.Output()
		if err != nil {
			fmt.Println("Error getting files of snapshot:", err)
			return
		}
	
		scanner := bufio.NewScanner(bytes.NewReader(output))
		if scanner.Scan() {
			// Skip the first line and continue scanning the rest
		}

		for scanner.Scan() {
			var entry SnapshotEntry
			err := json.Unmarshal([]byte(scanner.Text()), &entry)
			if err != nil {
				fmt.Println("Error parsing snapshot entry:", err)
				continue
			}
	
			accumulatedPaths = append(accumulatedPaths, entry.Path)
		}
	
		if err := scanner.Err(); err != nil {
			fmt.Println("Error scanning snapshot entries:", err)
			return
		}
	}
	
	// Print or use accumulatedPaths as needed
	uniquePaths := removeDuplicates(accumulatedPaths)
	sort.Strings(uniquePaths)
	fmt.Println("Accumulated Paths:")
	for _, path := range uniquePaths {
		fmt.Println(path)
	}
	
}

func removeDuplicates(input []string) []string {
	uniqueMap := make(map[string]struct{})
	uniqueSlice := []string{}

	for _, str := range input {
		if _, exists := uniqueMap[str]; !exists {
			uniqueMap[str] = struct{}{}
			uniqueSlice = append(uniqueSlice, str)
		}
	}

	return uniqueSlice
}