package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"
)

const (
	TEMP_SCRIPT_NAME = "tempscript_protobuf.go"
)

var (
	flagProtobufFileName = flag.String("p", "types.proto", "Specify path to proto-buf file.")
	flagJsonFileName     = flag.String("j", "types.json", "Specify path to json file that including content data.")
	flagPathUri          = flag.String("u", "http://localhost:8001", "Specify target uri.")
	flagRootType         = flag.String("r", "", "[MUST]Root type name")
)

type Template struct {
	RootType string
	JsonData string
	Endpoint string
}

func main() {
	flag.Parse()
	if *flagRootType == "" {
		flag.Usage()
		os.Exit(2)
	}

	err := wMain()
	if err != nil {
		os.Exit(2)
	}

	os.Exit(0)
}

func wMain() (err error) {
	err = CreateProtobuf("")
	if err != nil {
		return
	}

	code, err := CreateScript()
	if err != nil {
		return
	}

	err = WriteScript(code)
	defer func() { err = os.Remove(TEMP_SCRIPT_NAME) }()
	if err != nil {
		return
	}

	err = DoScript()
	return
}

func CreateProtobuf(filepath string) (err error) {
	cmd := exec.Command("protoc", "--go_out=.", *flagProtobufFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	return
}

func CreateScript() (code string, err error) {
	f, err := os.Open(*flagJsonFileName)
	if err != nil {
		return
	}

	json, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	t := template.Must(template.ParseFiles("runner.go.tmpl"))
	buf := bytes.NewBuffer([]byte{})
	t_dat := Template{
		RootType: *flagRootType,
		JsonData: string(json),
		Endpoint: *flagPathUri,
	}

	err = t.Execute(buf, t_dat)
	if err != nil {
		return
	}

	code = buf.String()
	return
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
	cmd := exec.Command("go", "run", TEMP_SCRIPT_NAME, "types.pb.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	return
}
