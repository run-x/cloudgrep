package config

import (
	"gopkg.in/yaml.v3"
)

var _ yaml.Unmarshaler = &NestedField{}

func (f NestedField) Empty() bool {
	return len(f) == 0
}

func (f NestedField) Last() Field {
	if len(f) == 0 {
		return Field{}
	}

	return f[len(f)-1]
}

func (f *NestedField) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.SequenceNode {
		return f.decodeSequenceNode(value)
	}
	if value.Kind == yaml.ScalarNode {
		return f.decodeScalarNode(value)
	}

	return &yaml.TypeError{Errors: []string{
		"unexpected node kind",
	}}
}

func (f *NestedField) decodeSequenceNode(value *yaml.Node) error {
	var fields []Field
	err := value.Decode(&fields)
	if err != nil {
		return err
	}

	*f = NestedField(fields)
	return nil
}

func (f *NestedField) decodeScalarNode(value *yaml.Node) error {
	var name string
	err := value.Decode(&name)
	if err != nil {
		return err
	}

	*f = NestedField{Field{Name: name}}
	return nil
}

func (f *Field) Zero() bool {
	if f == nil {
		return true
	}

	return f.Name == ""
}

var _ yaml.Unmarshaler = &Field{}

func (f *Field) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.MappingNode {
		return f.decodeMappingNode(value)
	} else if value.Kind == yaml.ScalarNode {
		return f.decodeScalarNode(value)
	}

	return &yaml.TypeError{Errors: []string{
		"unexpected node kind",
	}}
}

func (f *Field) decodeScalarNode(value *yaml.Node) error {
	var name string
	err := value.Decode(&name)
	if err != nil {
		return err
	}

	f.Name = name
	return nil
}

func (f *Field) decodeMappingNode(value *yaml.Node) error {
	// This must match the type def of Field
	// (but importantly, this type is not a yaml.Unmarshaler to avoid infinite recursion)
	field := struct {
		Name      string `yaml:"name"`
		SliceType string `yaml:"sliceType"`
		Pointer   bool   `yaml:"pointer"`
	}{}

	err := value.Decode(&field)
	if err != nil {
		return err
	}

	if field.Name == "" {
		return &yaml.TypeError{Errors: []string{
			"missing \"name\"",
		}}
	}

	*f = field
	return nil
}

func (t *TagField) Zero() bool {
	if t == nil {
		return true
	}

	return t.Field.Empty()
}

func TagFieldValidStyles() []string {
	return []string{
		"map",
		"struct",
	}
}
