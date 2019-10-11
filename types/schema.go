package types

import "github.com/hashicorp/terraform/helper/schema"

const (
	varsAttributeValues           = "values"
	varsAttributeJSON             = "json"
	inventoryAttributeHost        = "host"
	inventoryAttributeGroup       = "group"
	inventoryAttributeName        = "name"
	inventoryAttributeChildren    = "children"
	inventoryAttributeAlias       = "alias"
	inventoryAttributeAnsibleHost = "ansible_host"
	inventoryAttributeVariables   = "variables"
)

func varsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				varsAttributeValues: &schema.Schema{
					Type:     schema.TypeMap,
					Optional: true,
				},
				varsAttributeJSON: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func inventoryHostSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				inventoryAttributeAlias: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				inventoryAttributeAnsibleHost: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				inventoryAttributeVariables: varsSchema(),
			},
		},
	}
}

func inventoryGroupSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				inventoryAttributeName: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				inventoryAttributeHost: inventoryHostSchema(),
				inventoryAttributeVariables: varsSchema(),
			},
		},
	}
}
