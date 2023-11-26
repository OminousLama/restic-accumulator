package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	indexer "racc/modules/indexer"
	"racc/modules/util"

	"github.com/google/uuid"
)

var version = "undefined"
var metaBuildTime = "undefined"
var metaBuilderOS = "undefined"
var metaBuilderArch = "undefined"
var appname = "racc"

func main() {
	//#region Flags
	//#region Definition
	repoUrlPtr := flag.String("repoUrl", "", "Repository URL")
	repoSecretPtr := flag.String("repoSecret", "", "Repository Secret")
	dumpPathPrt := flag.String("dmpPath", "", "Directory a dump folder will be created in containing all node info in JSON format")
	//#endregion
	flag.Parse()
	//#region Validation
	if *repoUrlPtr == "" || *repoSecretPtr == "" {
		fmt.Printf("Usage: %s -repoUrl <url> -repoSecret <secret> -dmpPath <dmp_path>", appname)
		return
	}
	//#endregion
	//#endregion Flags

	allNodes := indexer.GetRepoFiles(*repoUrlPtr, *repoSecretPtr)
	onlyFiles := util.RemoveDirectories(allNodes)
	uniqueOnlyFiles := util.RemoveDuplicateFilenames(onlyFiles)
	for _, file := range uniqueOnlyFiles {
		fmt.Printf("File %s\n", file.Path)
	}

	if *dumpPathPrt == "" {
		return
	}

	dmpDirUUID := uuid.New().String()
	dmpDir := *dumpPathPrt + "/" + appname + "-dmp-" + dmpDirUUID + "/"
	util.EnsureDirExists(dmpDir)

	DumpToJson(dmpDir+"all-nodes.json", allNodes)
	DumpToJson(dmpDir+"file-nodes.json", onlyFiles)
	DumpToJson(dmpDir+"unique-file-nodes.json", uniqueOnlyFiles)
}

func DumpToJson(filepath string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Failed to write file:", err)
	}

	fmt.Printf("Dumped data to %s\n", filepath)
}
