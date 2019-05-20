package outgoingportal

import (
	"io/ioutil"
	"text/template"
)

var outgoingPortalEnvoyConfig *template.Template

func init() {
	envoyConfig, err := ioutil.ReadFile("config/outgoing-envoy.tpl")
	if err != nil {
		panic(err)
	}

	outgoingPortalEnvoyConfig = template.Must(template.New("config/outgoing-envoy.tpl").Parse(string(envoyConfig)))
}
