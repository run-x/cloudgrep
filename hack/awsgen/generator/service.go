package generator

import (
	"strings"

	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/template"
	"github.com/run-x/cloudgrep/hack/awsgen/util"
)

func (g Generator) generateService(service config.Service) string {
	var imports util.ImportSet

	buf := &strings.Builder{}

	reg, regImports := g.generateServiceRegister(service)
	buf.WriteString(reg)
	imports.Merge(regImports)

	for _, t := range service.Types {
		f, typeImports := g.generateType(service, t)
		buf.WriteString(f)
		imports.Merge(typeImports)
	}

	header := g.generateFileHeader(PackageName, imports.Get())

	return header + "\n" + buf.String()
}

func (g Generator) generateServiceRegister(service config.Service) (string, util.ImportSet) {
	data := struct {
		ProviderName string
		ServiceName  string
		FuncName     string

		Types []typeRegisterInfo
	}{
		ProviderName: ProviderStructName,
		ServiceName:  service.Name,
		FuncName:     registerFuncName(service),
	}

	var imports util.ImportSet

	for _, typ := range service.Types {
		data.Types = append(data.Types, typeRegisterInfo{
			ResourceName:  resourceName(service, typ),
			FetchFuncName: fetchFuncName(service, typ),
			IDField:       typ.ListAPI.IDField,
			Global:        typ.Global,
			TagField:      typ.ListAPI.Tags,
		})

		if !typ.GetTagsAPI.TagField.Zero() {
			imports.AddPath("github.com/run-x/cloudgrep/pkg/resourceconverter")
		}
	}

	return template.RenderTemplate("service-register.go", data), imports
}

type typeRegisterInfo struct {
	ResourceName  string
	FetchFuncName string

	IDField  config.Field
	Global   bool
	TagField config.TagField
}
