package types

import "github.com/fredwangwang/orderedmap"

type InventoryVariables map[string]interface{}

type InventoryHost struct {
	Alias         string
	AnsibleHost   string
	Variables     InventoryVariables
}

type InventoryGroup struct {
	Name          string
	Hosts         []InventoryHost
	Children      []InventoryGroup
	Variables     InventoryVariables
}

type InventoryRoot struct {
	Hosts         []InventoryHost
	Groups        []InventoryGroup
	Variables     InventoryVariables
}

type InventoryJSONGroup struct {
	Hosts    *orderedmap.OrderedMap `json:"hosts,omitempty"`
	Children InventoryJSONGroupList `json:"children,omitempty"`
	Vars     InventoryVariables     `json:"vars,omitempty"`
}

type InventoryJSONGroupList map[string]InventoryJSONGroup

type InventoryJSONRoot struct {
	All InventoryJSONGroup `json:"all,omitempty"`
}
