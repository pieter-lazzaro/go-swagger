// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
	"golang.org/x/tools/imports"
)

// Debug when the env var DEBUG is not empty
// the generators will be very noisy about what they are doing
var Debug = os.Getenv("DEBUG") != ""

var reservedGoWords = []string{
	"break", "default", "func", "interface", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
}

var defaultGoImports = []string{
	"bool", "int", "int8", "int16", "int32", "int64",
	"uint", "uint8", "uint16", "uint32", "uint64",
	"float32", "float64", "interface{}", "string",
	"byte", "rune",
}

var reservedGoWordSet map[string]struct{}

func init() {
	reservedGoWordSet = make(map[string]struct{})
	for _, gw := range reservedGoWords {
		reservedGoWordSet[gw] = struct{}{}
	}
}

func mangleName(name, suffix string) string {
	if _, ok := reservedGoWordSet[swag.ToFileName(name)]; !ok {
		return name
	}
	return strings.Join([]string{name, suffix}, "_")
}

func findSwaggerSpec(name string) (string, error) {
	f, err := os.Stat(name)
	if err != nil {
		return "", err
	}
	if f.IsDir() {
		return "", fmt.Errorf("%s is a directory", name)
	}
	return name, nil
}

// GenOpts the options for the generator
type GenOpts struct {
	Spec              string
	APIPackage        string
	ModelPackage      string
	ServerPackage     string
	ClientPackage     string
	Principal         string
	Target            string
	TypeMapping       map[string]string
	Imports           map[string]string
	DumpData          bool
	DefaultScheme     string
	IncludeModel      bool
	IncludeValidator  bool
	IncludeHandler    bool
	IncludeParameters bool
	IncludeResponses  bool
	IncludeMain       bool
	CustomFormatsFile string
	TemplateDir       string
}

type generatorOptions struct {
	ModelPackage    string
	TargetDirectory string
}

// on its way out
type propertyDescriptor struct {
	PropertyName      string
	ParamName         string
	Path              string
	ValueExpression   string
	IndexVar          string
	IsPrimitive       bool
	IsCustomFormatter bool
	IsContainer       bool
	IsMap             bool
}

// on its way out
type commonValidations struct {
	propertyDescriptor
	sharedValidations
	Type    string
	Format  string
	Items   *spec.Items
	Default interface{}
}

// on its way out
type genValidations struct {
	Type                string
	Required            bool
	DefaultValue        string
	MaxLength           int64
	MinLength           int64
	Pattern             string
	MultipleOf          float64
	Minimum             float64
	Maximum             float64
	ExclusiveMinimum    bool
	ExclusiveMaximum    bool
	Enum                string
	HasValidations      bool
	Format              string
	MinItems            int64
	MaxItems            int64
	UniqueItems         bool
	HasSliceValidations bool
	NeedsSize           bool
}

func loadSpec(specFile string) (string, *spec.Document, error) {
	// find swagger spec document, verify it exists
	specPath := specFile
	var err error
	if !strings.HasPrefix(specPath, "http") {
		specPath, err = findSwaggerSpec(specFile)
		if err != nil {
			return "", nil, err
		}
	}

	// load swagger spec
	specDoc, err := spec.Load(specPath)
	if err != nil {
		return "", nil, err
	}
	return specPath, specDoc, nil
}

func fileExists(target, name string) bool {
	ffn := swag.ToFileName(name) + ".go"
	_, err := os.Stat(filepath.Join(target, ffn))
	return !os.IsNotExist(err)
}

func writeToFileIfNotExist(target, name string, content []byte) error {
	if fileExists(target, name) {
		return nil
	}
	return writeToFile(target, name, content)
}

func formatGoFile(ffn string, content []byte) ([]byte, error) {
	opts := new(imports.Options)
	opts.TabIndent = true
	opts.TabWidth = 2
	opts.Fragment = true
	opts.Comments = true

	return imports.Process(ffn, content, opts)
}

func writeToFile(target, name string, content []byte) error {
	ffn := swag.ToFileName(name) + ".go"

	res, err := formatGoFile(ffn, content)
	if err != nil {
		log.Println(err)
		return writeFile(target, ffn, content)
	}

	return writeFile(target, ffn, res)
}

func writeFile(target, ffn string, content []byte) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(target, ffn), content, 0644)
}

func commentedLines(str string) string {
	lines := strings.Split(str, "\n")
	var commented []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			if !strings.HasPrefix(strings.TrimSpace(line), "//") {
				commented = append(commented, "// "+line)
			} else {
				commented = append(commented, line)
			}
		}
	}
	return strings.Join(commented, "\n")
}

func gatherModels(specDoc *spec.Document, modelNames []string) map[string]spec.Schema {
	models, mnc := make(map[string]spec.Schema), len(modelNames)
	for k, v := range specDoc.Spec().Definitions {
		if mnc == 0 {
			models[k] = v
		}
		for _, nm := range modelNames {
			if k == nm {
				models[k] = v
			}
		}
	}
	return models
}

func appNameOrDefault(specDoc *spec.Document, name, defaultName string) string {
	if name == "" {
		if specDoc.Spec().Info != nil && specDoc.Spec().Info.Title != "" {
			name = specDoc.Spec().Info.Title
		} else {
			name = defaultName
		}
	}
	return strings.TrimSuffix(swag.ToGoName(name), "API")
}

var namesCounter int64

func ensureUniqueName(key, method, path string, operations map[string]opRef) string {
	nm := key
	if nm == "" {
		nm = swag.ToGoName(strings.ToLower(method) + " " + path)
	}
	_, found := operations[nm]
	if found {
		namesCounter++
		return fmt.Sprintf("%s%d", nm, namesCounter)
	}
	return nm
}

func containsString(names []string, name string) bool {
	for _, nm := range names {
		if nm == name {
			return true
		}
	}
	return false
}

type opRef struct {
	Method string
	Path   string
	Op     spec.Operation
}

func gatherOperations(specDoc *spec.Document, operationIDs []string) map[string]opRef {
	operations := make(map[string]opRef)

	for method, pathItem := range specDoc.Operations() {
		for path, operation := range pathItem {
			if len(operationIDs) == 0 || containsString(operationIDs, operation.ID) {
				nm := ensureUniqueName(operation.ID, method, path, operations)
				vv := *operation
				vv.ID = nm
				operations[nm] = opRef{
					Method: method,
					Path:   path,
					Op:     vv,
				}
			}
		}
	}

	return operations
}

func pascalize(arg string) string {
	if len(arg) == 0 || arg[0] > '9' {
		return swag.ToGoName(arg)
	}

	return swag.ToGoName("Nr " + arg)
}
