// JSON file format info: https://help.salesforce.com/s/articleView?id=sf.cms_import_content_json.htm&type=5

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
)

type ContentBody struct {
	Title   string
	AltText string
}

type ContentItem struct {
	ContentType string
	UrlName     string
	ContentBody ContentBody
}

var contentItems []ContentItem

func isImageTypeSupported(fileName string) bool {
	matched, _ := regexp.MatchString("(\\S+(\\.(?i)(jpe?g|png|gif|bmp))$)", fileName)
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
		if isValidFile(file, cmsImportDir) {
			contentItem := ContentItem{
				ContentType: "cms_image",
				UrlName:     file.Name(),
				ContentBody: ContentBody{
					Title:   file.Name(),
					AltText: "alt text for " + file.Name(),
				},
			}
			contentItems = append(contentItems, contentItem)

		}
	}

	jsonOutput, err := json.MarshalIndent(contentItems, "", "\t")
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Println(string(jsonOutput))
	}
}
