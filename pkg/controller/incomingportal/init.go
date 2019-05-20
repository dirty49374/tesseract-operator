package incomingportal

import (
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
)

var tesseractNamespace = "sys-tesseract"
var tesseractHostname = ""
var tesseractReplicas int32 = 1

var incomingPortalEnvoyConfig *template.Template

func init() {
	if os.Getenv("TESSERACT_NAMESPACE") != "" {
		tesseractNamespace = os.Getenv("TESSERACT_NAMESPACE")
	}

	tesseractHostname = os.Getenv("TESSERACT_HOSTNAME")

	replicas, err := strconv.Atoi(os.Getenv("TESSERACT_REPLICAS"))
	if replicas != 0 || err != nil {
		tesseractReplicas = int32(replicas)
	}

	envoyConfig, err := ioutil.ReadFile("config/incoming-envoy.tpl")
	if err != nil {
		panic(err)
	}

	incomingPortalEnvoyConfig = template.Must(template.New("config/incoming-envoy.tpl").Parse(string(envoyConfig)))
}
