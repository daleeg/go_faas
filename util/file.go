package util

import (
	"fmt"
	"io/ioutil"
	"path"
)

func FindFile(directoryPath string) []string {
	baseFile, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		fmt.Println("An error occurred while open file :[" + directoryPath + "] .")
		fmt.Println(err)
		return nil
	}
	var res []string
	for _, fileItem := range baseFile {
		subPath := path.Join(directoryPath, fileItem.Name())
		if fileItem.IsDir() {
			innerFiles := FindFile(subPath)
			res = append(res, innerFiles...)
		} else {
			res = append(res, subPath)
		}
	}
	return res
}
