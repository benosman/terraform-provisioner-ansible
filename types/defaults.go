package types

import (
	"github.com/hashicorp/terraform/helper/schema"
)


// Defaults represents default settings for each consequent play.
type Defaults struct {
	hosts             []string
	groups            []string
	inventory         InventoryRoot
	becomeMethod      string
	becomeUser        string
	extraVars         []InventoryVariables
	forks             int
	inventoryFile     string
	limit             string
	vaultID           []string
	vaultPasswordFile string
	//
	hostsIsSet             bool
	groupsIsSet            bool
	inventoryIsSet         bool
	becomeMethodIsSet      bool
	becomeUserIsSet        bool
	extraVarsIsSet         bool
	forksIsSet             bool
	inventoryFileIsSet     bool
	limitIsSet             bool
	vaultIDIsSet           bool
	vaultPasswordFileIsSet bool
}

const (
	// attribute names:
	defaultsAttributeHosts             = "hosts"
	defaultsAttributeGroups            = "groups"
	defaultsAttributeBecomeMethod      = "become_method"
	defaultsAttributeBecomeUser        = "become_user"
	defaultsAttributeExtraVars         = "extra_vars"
	defaultsAttributeForks             = "forks"
	defaultsAttributeInventory         = "inventory"
	defaultsAttributeInventoryFile     = "inventory_file"
	defaultsAttributeLimit             = "limit"
	defaultsAttributeVaultID           = "vault_id"
	defaultsAttributeVaultPasswordFile = "vault_password_file"
	defaultsAttributeHostsAlias		   = "alias"
	defaultsAttributeHostsHost         = "ansible_host"
	defaultsAttributeHostsVariables    = "variables"
	defaultsAttributeHostsVariablesJSON = "variables_json"

)

// NewDefaultsSchema returns a new defaults schema.
func NewDefaultsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				defaultsAttributeHosts: &schema.Schema{
					Type:     schema.TypeList,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Optional: true,
				},
				defaultsAttributeGroups: &schema.Schema{
					Type:     schema.TypeList,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Optional: true,
				},
				defaultsAttributeInventory: &schema.Schema{
					Type:     schema.TypeList,
					Optional: true,
					MaxItems:    1,
					Elem:     &schema.Resource{
						Schema: map[string]*schema.Schema{
							inventoryAttributeHost:  inventoryHostSchema(),
							inventoryAttributeGroup: inventoryGroupSchema(),
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
					ConflictsWith: []string{"defaults.hosts", "defaults.groups"},
				},
				defaultsAttributeBecomeMethod: &schema.Schema{
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: vfBecomeMethod,
				},
				defaultsAttributeBecomeUser: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				defaultsAttributeExtraVars: extraVarsSchema(),
				defaultsAttributeForks: &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
				},
				defaultsAttributeInventoryFile: &schema.Schema{
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: vfPath,
				},
				defaultsAttributeLimit: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				defaultsAttributeVaultID: &schema.Schema{
					Type:          schema.TypeList,
					Elem:          &schema.Schema{Type: schema.TypeString},
					Optional:      true,
					ConflictsWith: []string{"defaults.vault_password_file"},
				},
				defaultsAttributeVaultPasswordFile: &schema.Schema{
					Type:          schema.TypeString,
					Optional:      true,
					ValidateFunc:  vfPath,
					ConflictsWith: []string{"defaults.vault_id"},
				},
			},
		},
	}
}

// NewDefaultsFromInterface reads Defaults configuration from Terraform schema.
func NewDefaultsFromInterface(i interface{}, ok bool) (*Defaults, error) {
	if ok {
		vals := mapFromTypeSetList(i.(*schema.Set).List())
		return NewDefaultsFromMapInterface(vals, ok)
	}
	return &Defaults{}, nil
}

// NewDefaultsFromMapInterface reads Defaults configuration from a map.
func NewDefaultsFromMapInterface(vals map[string]interface{}, ok bool) (*Defaults, error) {
	v := &Defaults{}
	if ok {
		if val, ok := vals[defaultsAttributeHosts]; ok {
			v.hosts = listOfInterfaceToListOfString(val.([]interface{}))
			v.hostsIsSet = len(v.hosts) > 0
		}
		if val, ok := vals[defaultsAttributeGroups]; ok {
			v.groups = listOfInterfaceToListOfString(val.([]interface{}))
			v.groupsIsSet = len(v.groups) > 0
		}
		if val, ok := vals[defaultsAttributeInventory]; ok {
			inventory, err := listOfInterfaceToInventoryRoot(val.([]interface{}), "defaults")
			if err != nil {
				return nil, err
			}
			v.inventory = inventory
			v.inventoryIsSet = len(v.inventory.Hosts) > 0
		}
		if val, ok := vals[defaultsAttributeBecomeMethod]; ok {
			v.becomeMethod = val.(string)
			v.becomeMethodIsSet = v.becomeMethod != ""
		}
		if val, ok := vals[defaultsAttributeBecomeUser]; ok {
			v.becomeUser = val.(string)
			v.becomeUserIsSet = v.becomeUser != ""
		}
		if val, ok := vals[defaultsAttributeExtraVars]; ok && len(val.([]interface{})) > 0 {
			extraVars, err := listOfInterfaceToExtraVars(val, "defaults")
			if err != nil {
				return nil, err
			}
			v.extraVars = extraVars
			v.extraVarsIsSet = len(v.extraVars) > 0
		}
		if val, ok := vals[defaultsAttributeForks]; ok {
			v.forks = val.(int)
			v.forksIsSet = v.forks > 0
		}
		if val, ok := vals[defaultsAttributeInventoryFile]; ok {
			v.inventoryFile = val.(string)
			v.inventoryFileIsSet = v.inventoryFile != ""
		}
		if val, ok := vals[defaultsAttributeLimit]; ok {
			v.limit = val.(string)
			v.limitIsSet = v.limit != ""
		}
		if val, ok := vals[defaultsAttributeVaultID]; ok {
			v.vaultID = listOfInterfaceToListOfString(val.([]interface{}))
			v.vaultIDIsSet = len(v.vaultID) > 0
		}
		if val, ok := vals[defaultsAttributeVaultPasswordFile]; ok {
			v.vaultPasswordFile = val.(string)
			v.vaultPasswordFileIsSet = v.vaultPasswordFile != ""
		}
	}
	return v, nil
}

// Hosts returns default hosts.
func (v *Defaults) Hosts() []InventoryHost {
	//return v.hosts
	return nil
}

// BecomeMethod returns become method.
func (v *Defaults) BecomeMethod() string {
	return v.becomeMethod
}

// BecomeUser returns become user.
func (v *Defaults) BecomeUser() string {
	return v.becomeUser
}
