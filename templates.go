package main

import "fmt"

type projectTemplate struct {
	RequiredVars map[string]bool
}

var availableTemplates = map[string]*projectTemplate{
	"go-daemon": {
		RequiredVars: map[string]bool{
			"ModuleName":     true,
			"GoVersion":      true,
			"AppName":        true,
			"AppDescription": true,
		},
	},
}

func validateVarsMap(templateName string, varsMap map[string]string) bool {
	requiredVars := availableTemplates[templateName].RequiredVars
	is_valid := true
	for k := range requiredVars {
		if varsMap[k] == "" {
			is_valid = false
			fmt.Printf("Required %s variable but not provided\n", k)
		}
	}
	return is_valid
}
