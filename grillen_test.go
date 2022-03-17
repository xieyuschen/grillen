package grillen_test

import (
	"os"
	"testing"
	"text/template"
)

func TestSimpleTemplateGenerate(t *testing.T) {
	type Inventory struct {
		Material string
		Count    uint
	}
	sweaters := Inventory{"wool", 17}
	tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, sweaters)
	if err != nil {
		panic(err)
	}
}

const EchoTemplate = `
package main
import (
	"fmt"
)
func main(){
	fmt.Println("{{.EchoName}}")
}
`
const generatedFolder = "generated-file/"

type Root struct {
	EchoName string
}

func TestGenerateEchoFunction(t *testing.T) {
	data := Root{EchoName: "this is xieyuschen"}
	f, err := os.OpenFile(generatedFolder+"echo.go", os.O_CREATE|os.O_RDWR, 0755)
	tmpl, err := template.New("test").Parse(EchoTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, data)
	if err != nil {
		panic(err)
	}
}
