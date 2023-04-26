// JSON file format info:
// - https://help.salesforce.com/s/articleView?id=sf.cms_import_content_json.htm&type=5
// - https://help.salesforce.com/s/articleView?id=sf.cms_import_export_overview.htm&type=5

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ContentSource struct {
	Source string `json:"ref"`
}

type ContentBody struct {
	Title         string        `json:"title"`
	AltText       string        `json:"altText"`
	ContentSource ContentSource `json:"source"`
}

type ContentItem struct {
	ContentType string      `json:"type"`
	UrlName     string      `json:"urlName"`
	ContentBody ContentBody `json:"body"`
}

type Content struct {
	Content []ContentItem `json:"content"`
}

var contentItems []ContentItem

func reportErrorStdOut(errMsg string) {
	_, err := os.Stderr.WriteString(errMsg)
	if err != nil {
		panic(err)
	}
}

func isImageTypeSupported(fileName string) bool {
	matched, _ := regexp.MatchString("(\\S+(\\.(?i)(jpe?g|png|gif|bmp))$)", fileName)
	return matched
}

func isValidFile(file os.DirEntry, cmsImportDir string) bool {
	if file.IsDir() {
		reportErrorStdOut(file.Name() + " is a directory, skipping\n")
		return false
	}

	// per SFDC docs, image files cannot be >25MB
	fileInfo, _ := os.Lstat(cmsImportDir + "/_media/" + file.Name())
	if fileInfo.Size() > 2.5e+7 {
		reportErrorStdOut(file.Name() + " is greater than 25MB, skipping\n")
		return false
	}

	if !isImageTypeSupported(file.Name()) {
		reportErrorStdOut(file.Name() + " is of a non-supported file type, skipping\n")
		return false
	}
	return true
}

func main() {
	args := os.Args
	cmsImportDir := args[1]
	jsonFilename := args[2]
	files, err := os.ReadDir(cmsImportDir + "/_media")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if isValidFile(file, cmsImportDir) {
			fileNameWithoutExtension := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			contentItem := ContentItem{
				ContentType: "cms_image",
				UrlName:     fileNameWithoutExtension,
				ContentBody: ContentBody{
					Title:         fileNameWithoutExtension,
					AltText:       "alt text for " + fileNameWithoutExtension,
					ContentSource: ContentSource{file.Name()},
				},
			}
			contentItems = append(contentItems, contentItem)

		}
	}

	contentContainer := Content{Content: contentItems}
	jsonOutput, err := json.MarshalIndent(contentContainer, "", "   ")
	if err != nil {
		reportErrorStdOut(err.Error())
	} else {
		err = os.WriteFile(cmsImportDir+"/"+jsonFilename, jsonOutput, 0644)
	}
}
