package render

import (
	"fmt"
	"github.com/Galagoshin/GoLogger/logger"
	"github.com/Galagoshin/GoUtils/files"
	"io/ioutil"
	"os"
	"path/filepath"
)

var templates_storage map[string]GHtml

func appendFiles(arr []string, path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Panic(err)
	}

	for _, dir := range files {
		if dir.IsDir() {
			arr = appendFiles(arr, path+dir.Name()+"/")
		}
	}
	to_append, err := filepath.Glob(path + "*.html")
	if err != nil {
		logger.Panic(err)
	}
	arr = append(arr, to_append...)
	return arr
}

func LoadTemplates() {
	templates_storage = make(map[string]GHtml)
	viewsDir := files.Directory{Path: "views"}
	err := viewsDir.CreateAll()
	if err != nil {
		logger.Panic(err)
	}

	all_files := appendFiles([]string{}, "views/")
	for _, filename := range all_files {
		file := files.File{Path: filename}
		err := file.Open(os.O_RDWR)
		if err != nil {
			logger.Panic(err)
		}
		content := GHtml(file.ReadString())
		err = file.Close()
		if err != nil {
			logger.Panic(err)
		}
		templates_storage[filename[6:len(filename)-5]] = content
		logger.Debug(2, false, fmt.Sprintf("Registered template: %s (%s)", filename, filename[6:len(filename)-5]))
	}
	if len(all_files) > 0 {
		logger.Print(fmt.Sprintf("Loaded %d templates.", len(all_files)))
	}
	logger.Debug(10, false, fmt.Sprintf("Templates: %+v", templates_storage))
}

func GetView(template string) GHtml {
	return templates_storage[template]
}
