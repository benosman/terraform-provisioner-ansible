package types

import "github.com/hashicorp/terraform/helper/schema"

const (
	inventoryAttributeHosts 	    = "hosts"
	inventoryAttributeGroups 	    = "groups"
	inventoryAttributeName			= "name"
	inventoryAttributeChildren 	    = "children"
	inventoryAttributeAlias 	    = "alias"
	inventoryAttributeAnsibleHost   = "ansible_host"
	inventoryAttributeVariables     = "variables"
	inventoryAttributeVariablesJSON = "variables_json"
)


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
				inventoryAttributeVariables: &schema.Schema{
					Type:     schema.TypeMap,
					Optional: true,
				},
				inventoryAttributeVariablesJSON: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
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
				inventoryAttributeHosts: inventoryHostSchema(),
				inventoryAttributeVariables: &schema.Schema{
					Type:     schema.TypeMap,
					Optional: true,
				},
				inventoryAttributeVariablesJSON: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}