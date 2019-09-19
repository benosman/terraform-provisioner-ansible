package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/radekg/terraform-provisioner-ansible/ansible"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// ExportSchema should be called to export the structure
// of the provisioner.
func Export(p *schema.Provisioner) *ResourceProvisionerSchema {
	result := new(ResourceProvisionerSchema)

	result.Name = "ansible"
	result.Type = "provisioner"
	result.Version = "v2.3.0"
	result.Schema = schemaMap(p.Schema).Export()

	return result
}

func ExportResource(r *schema.Resource) SchemaInfo {
	return schemaMap(r.Schema).Export()
}

// schemaMap is a wrapper that adds nice functions on top of schemas.
type schemaMap map[string]*schema.Schema

// Export exports the format of this schema.
func (m schemaMap) Export() SchemaInfo {
	result := make(SchemaInfo)
	for k, v := range m {
		item := export(v)
		result[k] = item
	}
	return result
}

func export(v *schema.Schema) SchemaDefinition {
	item := SchemaDefinition{}

	item.Type = fmt.Sprintf("%s", v.Type)
	item.Optional = v.Optional
	item.Required = v.Required
	item.Description = v.Description
	item.InputDefault = v.InputDefault
	item.Computed = v.Computed
	item.MaxItems = v.MaxItems
	item.MinItems = v.MinItems
	item.PromoteSingle = v.PromoteSingle
	item.ComputedWhen = v.ComputedWhen
	item.ConflictsWith = v.ConflictsWith
	item.Deprecated = v.Deprecated
	item.Removed = v.Removed

	if v.Elem != nil {
		item.Elem = exportValue(v.Elem, fmt.Sprintf("%T", v.Elem))
	}

	// TODO: Find better solution
	if defValue, err := v.DefaultValue(); err == nil && defValue != nil && !reflect.DeepEqual(defValue, v.Default) {
		item.Default = exportValue(defValue, fmt.Sprintf("%T", defValue))
	}
	return item
}

func exportValue(value interface{}, t string) SchemaElement {
	s2, ok := value.(*schema.Schema)
	if ok {
		return SchemaElement{Type: "SchemaElements", ElementsType: fmt.Sprintf("%s", s2.Type)}
	}
	r2, ok := value.(*schema.Resource)
	if ok {
		return SchemaElement{Type: "SchemaInfo", Info: ExportResource(r2)}
	}
	return SchemaElement{Type: t, Value: fmt.Sprintf("%v", value)}
}

func Generate(provisioner *schema.Provisioner, name string, outputPath string) {
	outputFilePath := filepath.Join(outputPath, fmt.Sprintf("%s.json", name))

	if err := DoGenerate(provisioner, name, outputFilePath); err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		os.Exit(255)
	}
}

func DoGenerate(provisioner *schema.Provisioner, provisionerName string, outputFilePath string) error {
	provisionerJson, err := json.MarshalIndent(Export(provisioner), "", "  ")

	if err != nil {
		return err
	}

	file, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(provisionerJson)
	if err != nil {
		return err
	}

	return file.Sync()
}

type SchemaElement struct {
	// One of ValueType or "SchemaElements" or "SchemaInfo"
	Type string `json:",omitempty"`
	// Set for simple types (from ValueType)
	Value string `json:",omitempty"`
	// Set if Type == "SchemaElements"
	ElementsType string `json:",omitempty"`
	// Set if Type == "SchemaInfo"
	Info SchemaInfo `json:",omitempty"`
}

type SchemaDefinition struct {
	Type          string `json:",omitempty"`
	Optional      bool   `json:",omitempty"`
	Required      bool   `json:",omitempty"`
	Description   string `json:",omitempty"`
	InputDefault  string `json:",omitempty"`
	Computed      bool   `json:",omitempty"`
	MaxItems      int    `json:",omitempty"`
	MinItems      int    `json:",omitempty"`
	PromoteSingle bool   `json:",omitempty"`

	ComputedWhen  []string `json:",omitempty"`
	ConflictsWith []string `json:",omitempty"`

	Deprecated string `json:",omitempty"`
	Removed    string `json:",omitempty"`

	Default SchemaElement `json:",omitempty"`
	Elem    SchemaElement `json:",omitempty"`
}

type SchemaInfo map[string]SchemaDefinition

// ResourceProvisionerSchema
type ResourceProvisionerSchema struct {
	Name        string                `json:"name"`
	Type        string                `json:"type"`
	Version     string                `json:"version"`
	Schema      SchemaInfo            `json:"schema"`
}

func main() {
	var provisioner terraform.ResourceProvisioner
	provisioner = ansible.Provisioner()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	Generate(provisioner.(*schema.Provisioner), "ansible-provisioner",  dir)
}
