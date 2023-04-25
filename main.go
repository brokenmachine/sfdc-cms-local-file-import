// CMS import info: https://help.salesforce.com/s/articleView?id=sf.cms_customcontenttypes.htm&type=5

package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func isImageTypeSupported(fileName string) bool {
	matched, _ := regexp.MatchString("([^\\s]+(\\.(?i)(jpe?g|png|gif|bmp))$)", fileName)
	return matched
}

func isValidFile(file os.DirEntry, cmsImportDir string) bool {
	if file.IsDir() {
		os.Stderr.WriteString(file.Name() + " is a directory, skipping\n")
		return false
	}

	// per SFDC docs, image files cannot be >25MB
	fileInfo, _ := os.Lstat(cmsImportDir + "/_media/" + file.Name())
	if fileInfo.Size() > 2.5e+7 {
		os.Stderr.WriteString(file.Name() + " is greater than 25MB, skipping\n")
		return false
	}

	if !isImageTypeSupported(file.Name()) {
		os.Stderr.WriteString(file.Name() + " is of a non-supported file type, skipping\n")
		return false
	}
	return true
}

func main() {
	args := os.Args
	cmsImportDir := args[1]

	files, err := os.ReadDir(cmsImportDir + "/_media")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name() + "," + strconv.FormatBool(isValidFile(file, cmsImportDir)))
	}
}
