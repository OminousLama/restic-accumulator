package indexer

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	types "racc/modules/types"
	"strings"
)

func GetRepoFiles(repoUrl string, repoSecret string) []types.File {
	snapshots := GetRepoSnapshots(repoUrl, repoSecret)
	var files []types.File
	for _, snapshot := range snapshots {
		cmd := exec.Command("restic", "ls", "--long", "-r", repoUrl, "--json", snapshot.ID, "--password-command", fmt.Sprintf("echo %s", repoSecret))

		output := RunCommandGetFullOutput(*cmd)

		preprocessedOutput := "[" + strings.ReplaceAll(string(output), "}{", "},{") + "]"

		var snapFiles []types.File
		err := json.NewDecoder(strings.NewReader(preprocessedOutput)).Decode(&snapFiles)
		if err != nil {
			fmt.Printf("Error decoding files of snapshot %s: %s", snapshot.ID, err)
			return nil
		}

		files = append(files, snapFiles...)
	}

	return files
}

func GetRepoSnapshots(repoUrl string, repoSecret string) []types.Snapshot {
	cmd := exec.Command("restic", "snapshots", "-r", repoUrl, "--json", "--password-command", fmt.Sprintf("echo %s", repoSecret))

	preprocessedOutput := RunCommandGetFullOutput(*cmd)

	var snapshots []types.Snapshot
	err := json.NewDecoder(strings.NewReader(preprocessedOutput)).Decode(&snapshots)
	if err != nil {
		fmt.Println("Error decoding snapshot JSON:", err)
		return nil
	}

	return snapshots
}

func RunCommandGetFullOutput(cmd exec.Cmd) string {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return ""
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return ""
	}

	output, err := io.ReadAll(stdout)
	if err != nil {
		fmt.Println("Error reading stdout:", err)
		return ""
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Command failed:", err)
		return ""
	}

	return strings.ReplaceAll(string(output), "\n", "")
}
