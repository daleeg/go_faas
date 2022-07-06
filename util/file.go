package util

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
)

var (
	logger *logrus.Logger
)

func FindFile(directoryPath string) []string {
	baseFile, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		logger.Errorln("An error occurred while open file :[" + directoryPath + "] .")
		logger.Errorln(err)
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

func init() {
	logger = GetLogger("util")
}
