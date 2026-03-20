package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

//go:embed templates
var templatesFS embed.FS

var version = "debug"

func main() {
	var vars Vars
	templateName := flag.String("template", "go-daemon", "The name of the template to use")
	versionFlag := flag.Bool("version", false, "Display current version")
	flag.Var(&vars, "var", "key=value pairs")
	flag.Parse()
	if *versionFlag {
		log.Println("initium", version)
		return
	}
	if availableTemplates[*templateName] == nil {
		fmt.Printf("No project template named '%s'\n", *templateName)
		return
	}
	varsMap := parseVars(vars)
	if !validateVarsMap(*templateName, varsMap) {
		return
	}
	templatePath := "templates/" + *templateName
	fs.WalkDir(templatesFS, templatePath, func(entryPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relativePath := strings.TrimPrefix(entryPath, templatePath+"/")
		if relativePath == templatePath {
			return nil
		}
		targetPath := path.Join(".", relativePath)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}
		content, err := templatesFS.ReadFile(entryPath)
		if err != nil {
			return err
		}
		targetPath = strings.TrimSuffix(targetPath, ".tmpl")
		tmpl, err := template.New(relativePath).Parse(string(content))
		if err != nil {
			return err
		}
		f, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer f.Close()
		err = tmpl.Execute(f, varsMap)
		if err != nil {
			return err
		}
		return nil
	})
}
