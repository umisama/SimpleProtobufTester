package main

import (
	"bytes"
	"os"
	"os/exec"
	"text/template"
)

const (
	TEMP_SCRIPT_NAME = "tempscript_protobuf.go"
)

type Template struct {
	Echo string
}

func main() {
	err := CreateProtobuf("")
	if err != nil {
		os.Exit(2)
	}

	code, err := CreateScript()
	if err != nil {
		os.Exit(2)
	}

	err = WriteScript(code)
	if err != nil {
		os.Exit(2)
	}
	println(code)

	err = DoScript()
	if err != nil {
		os.Exit(2)
	}
}

func CreateProtobuf(filepath string) error {
	return nil
}

func CreateScript() (code string, err error) {
	t := template.Must(template.New("script").Parse(code_template))

	buf := bytes.NewBuffer([]byte{})

	dat := Template{Echo: "hellohello from code"}

	t.Execute(buf, dat)
	return buf.String(), nil
}

func WriteScript(code string) (err error) {
	f, err := os.OpenFile(TEMP_SCRIPT_NAME, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}

	_, err = f.WriteString(code)
	return
}

func DoScript() (err error) {
	cmd := exec.Command("go", "run", TEMP_SCRIPT_NAME)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return
	}

	err = os.Remove(TEMP_SCRIPT_NAME)
	return
}

var code_template = `package main

import (
	"fmt"
)

func main() {
	fmt.Printf("{{.Echo}}")
}`
