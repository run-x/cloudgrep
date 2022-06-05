package generator

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/run-x/cloudgrep/hack/awsgen/config"
	"github.com/run-x/cloudgrep/hack/awsgen/util"
	"github.com/run-x/cloudgrep/hack/awsgen/writer"
)

type Generator struct {
	Format      bool
	LineNumbers bool
}

func (g *Generator) Generate(w writer.Writer, cfg config.Config) error {
	for _, service := range cfg.Services {
		text := g.generateService(service)
		err := g.writeFile(w, service.Name, text)
		if err != nil {
			return fmt.Errorf("cannot generate %s: %w", service.Name, err)
		}
	}

	text := g.generateMainFile(cfg.Services)
	err := g.writeFile(w, "services", text)
	if err != nil {
		return fmt.Errorf("cannot generate all services file: %w", err)
	}

	return nil
}

func (g *Generator) writeFile(w writer.Writer, name, text string) error {
	if g.Format {
		formatted, err := format.Source([]byte(text))
		if err != nil {
			return fmt.Errorf("cannot format text: %w", err)
		}

		text = string(formatted)
	}

	if g.LineNumbers {
		text = linenumbers(text)
	}

	err := w.WriteFile(name, []byte(text))
	if err != nil {
		return fmt.Errorf("cannot write %s: %w", name, err)
	}

	return nil
}

func (g *Generator) generateService(service config.ServiceConfig) string {
	buf := &strings.Builder{}

	buf.WriteString(g.generateServiceHeader(service))

	for _, t := range service.Types {
		c := listFuncConfig{
			ResourceName: fmt.Sprintf("%s.%s", service.Name, t.Name),
			FuncName:     fetchFuncName(service, t),
			Action:       t.ListAPI.Call,
			Paginated:    t.ListAPI.Pagination,
			Description:  t.Description,
			Type:         t.Name,
			ServicePkg:   service.Name,
			ProviderName: "Provider",
			OutputKey: util.RecursiveAppend{
				Keys: t.ListAPI.OutputKey,
				Root: "output",
			},
		}

		f := g.generateListFunction(c)
		buf.WriteString(f)

		if t.GetTagsAPI.Call != "" {
			f = g.generateTypeTagFunction(service, t)
			buf.WriteString(f)
		}
	}

	return buf.String()
}

func fetchFuncName(svc config.ServiceConfig, typ config.TypeConfig) string {
	return fmt.Sprintf(
		"fetch_%s_%s",
		svc.Name,
		typ.Name,
	)
}