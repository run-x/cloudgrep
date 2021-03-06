package generator

import (
	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/template"
)

// generateRegisterFile defines the provider-wide file that calls each service's registration function.
func (g Generator) generateRegisterFile(services []config.Service) string {
	data := struct {
		ProviderName      string
		RegisterFuncNames []string
	}{
		ProviderName: ProviderStructName,
	}

	for _, service := range services {
		data.RegisterFuncNames = append(data.RegisterFuncNames, registerFuncName(service))
	}

	return g.generateFileHeader(PackageName, nil) + template.RenderTemplate("register.go", data)
}
