package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"path"

	"os"
	"strings"
	"text/template"
)

//go:embed all:templates/*
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
		trimmedPath := strings.TrimPrefix(entryPath, templatePath)
		trimmedPath = strings.TrimPrefix(trimmedPath, "/")
		if trimmedPath == "" {
			return nil
		}
		targetPath := path.Join(".", trimmedPath)
		if d.IsDir() {
			err := os.MkdirAll(targetPath, 0755)
			if err != nil {
				fmt.Println(err)
			}
			return nil
		}
		content, err := templatesFS.ReadFile(entryPath)
		if err != nil {
			return err
		}
		targetPath = strings.TrimSuffix(targetPath, ".tmpl")
		err = os.MkdirAll(path.Dir(targetPath), 0755)
		if err != nil {
			fmt.Println(err)
			return err
		}
		tmpl, err := template.New(trimmedPath).Parse(string(content))
		if err != nil {
			fmt.Println(err)
			return err
		}
		f, err := os.Create(targetPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer f.Close()
		err = tmpl.Execute(f, varsMap)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
}
