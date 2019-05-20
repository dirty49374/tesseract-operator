package util

import (
	"bytes"
	"fmt"
	"text/template"
)

func ExpandTemplate(tpl string, data interface{}) (string, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println(err)
		}
	}()

	var buf bytes.Buffer

	template, err := template.New("envoyConfig").Parse(tpl)
	if err != nil {
		return "", err
	}

	err = template.Execute(&buf, data)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	return buf.String(), nil
}
