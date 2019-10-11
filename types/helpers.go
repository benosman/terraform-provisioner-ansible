package types

import (
	"encoding/json"
	"fmt"
	"github.com/fredwangwang/orderedmap"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

var (
	becomeMethods = map[string]bool{
		"sudo":   true,
		"su":     true,
		"pbrun":  true,
		"pfexec": true,
		"doas":   true,
		"dzdo":   true,
		"ksu":    true,
		"runas":  true,
	}
)

func vfBecomeMethod(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if !becomeMethods[v] {
		errs = append(errs, fmt.Errorf("%s is not a valid become_method", v))
	}
	return
}

func vfPath(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if strings.Index(v, "${path.module}") > -1 {
		warns = append(warns, fmt.Sprintf("Unable to determine the existence of '%s', most likely because of https://github.com/hashicorp/terraform/issues/17439. If the file does not exist, you'll experience a failure at runtime.", v))
	} else {
		if _, err := ResolvePath(v); err != nil {
			errs = append(errs, fmt.Errorf("file '%s' does not exist", v))
		}
	}
	return
}

// VfPathDirectory validates existence of a path and that the path is a directory.
func VfPathDirectory(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if strings.Index(v, "${path.module}") > -1 {
		warns = append(warns, fmt.Sprintf("Unable to determine the existence of '%s', most likely because of https://github.com/hashicorp/terraform/issues/17439. If the file does not exist, you'll experience a failure at runtime.", v))
	} else {
		if _, err := ResolveDirectory(v); err != nil {
			errs = append(errs, fmt.Errorf("directory '%s' does not exist or path is not a directory", v))
		}
	}
	return
}

func mapFromTypeMap(v interface{}) map[string]interface{} {
	switch v := v.(type) {
	case nil:
		return make(map[string]interface{})
	case map[string]interface{}:
		return v
	default:
		panic(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func listOfMapFromTypeMap(v interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	switch v := v.(type) {
	case nil:
		//do nothing
	case map[string]interface{}:
		if len(v) > 0 {
			result = append(result, mapFromTypeMap(v))
		}
	default:
		panic(fmt.Sprintf("Unsupported type: %T", v))
	}
	return result
}

func mapFromTypeSet(i interface{}) map[string]interface{} {
	return i.(map[string]interface{})
}

func mapFromTypeSetList(i []interface{}) map[string]interface{} {
	for _, v := range i {
		return mapFromTypeSet(v)
	}
	return make(map[string]interface{})
}

func listOfInterfaceToListOfString(v interface{}) []string {
	var result []string
	switch v := v.(type) {
	case nil:
		return result
	case []interface{}:
		for _, vv := range v {
			if vv, ok := vv.(string); ok {
				result = append(result, vv)
			}
		}
		return result
	default:
		panic(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func listOfInterfaceToListOfHosts(v interface{}) ([]InventoryHost, error)  {
	var result = make([]InventoryHost, 0)
	switch v := v.(type) {
	case nil:
		return nil, nil
	case []interface{}:
		for _, vv := range v {
			hostMap := vv.(map[string]interface {})
			host := InventoryHost{}
			if value, found := hostMap[inventoryAttributeAlias]; found {
				host.Alias =  value.(string);
			}
			if value, found := hostMap[inventoryAttributeAnsibleHost]; found {
				host.AnsibleHost =  value.(string);
			}
			if value, found := hostMap[inventoryAttributeVariables]; found {
				host.Variables = value.(map[string]interface {});
			}
			result = append(result, host)
		}
		return result, nil
	default:
		return result, fmt.Errorf("Unsupported type: %T", v)
	}
}

func listOfStringToListOfHosts(v [] string) []InventoryHost {
	var result = make([]InventoryHost, 0)
	for _, vv := range v {
		host := InventoryHost{}
		host.Alias = vv
		result = append(result, host)
	}
	return result
}

func listOfInterfaceToListOfGroups(v interface{}) ([]InventoryGroup, error) {
	var result = make([]InventoryGroup, 0)
	switch v := v.(type) {
	case nil:
		return nil, nil
	case []interface{}:
		for _, vv := range v {
			groupMap := vv.(map[string]interface {})
			group := InventoryGroup{}
			if value, found := groupMap[inventoryAttributeName]; found {
				group.Name =  value.(string)
			}
			if value, found := groupMap[inventoryAttributeHost]; found {
				hosts, err :=  listOfInterfaceToListOfHosts(value.([]interface{}))
				if err != nil {
					return result, err
				}
				group.Hosts = hosts
			}
			if value, found := groupMap[inventoryAttributeVariables]; found {
				group.Variables = value.(map[string]interface {})
			}
			result = append(result, group)
		}
		return result, nil
	default:
		return result, fmt.Errorf("Unsupported type: %T", v)
	}
}

func listOfStringToListOfGroups(v [] string, hosts []InventoryHost) ([]InventoryGroup, error)  {
	var result = make([]InventoryGroup, 0)
	for _, vv := range v {
		group := InventoryGroup{}
		group.Name = vv
		group.Hosts = hosts
		result = append(result, group)
	}
	return result, nil
}

func listOfInterfaceToExtraVars(v interface{}, key string) ([]InventoryVariables, error) {
	var result = make([]InventoryVariables, 0)
	switch v := v.(type) {
	case nil:
		return result, nil

	case []interface{}:
		for index, vv := range v {

			if item, ok := vv.(map[string]interface{}); ok {
				vars := InventoryVariables{}

				if value, found := item[extraVarsAttributeValues]; found {
					vars = value.(map[string]interface {})
				}
				if value, found := item[extraVarsAttributeJSON]; found {
					json := value.(string)
					if len(json) > 0 {
						if len(vars) > 0 {
							return result, fmt.Errorf("%s.extra_vars[%d]: invalid syntax: both values and json attributes are set", key, index)
						}
						jsonVars, err := variablesJSONToVariablesMap(json)
						if err != nil {
							return result, fmt.Errorf("%s.extra_vars[%d]: %s", key, index, err)
						}
						vars = jsonVars
					}
				}
				result = append(result, vars)
			}
		}
		return result, nil
	default:
		return result, fmt.Errorf("Unsupported type: %T", v)
	}
}


func listOfInterfaceToInventoryRoot(v interface{}, key string) (InventoryRoot, error) {
	var root = InventoryRoot{}
	switch v := v.(type) {
	case nil:
		return root, nil

	case []interface{}:
		if len(v) > 0 {
			inventoryMap := v[0].(map[string]interface{})
			if value, found := inventoryMap[inventoryAttributeHost]; found {
				hosts, err := listOfInterfaceToListOfHosts(value.([]interface{}))
				if err != nil {
					return root, err
				}
				root.Hosts = hosts
			}
			if value, found := inventoryMap[inventoryAttributeGroup]; found {
				groups, err :=listOfInterfaceToListOfGroups(value.([]interface{}))
				if err != nil {
					return root, err
				}
				root.Groups = groups
			}
			if value, found := inventoryMap[inventoryAttributeVariables]; found {
				root.Variables =value.(map[string]interface{})
			}
			if value, found := inventoryMap[inventoryAttributeVariablesJSON]; found {
				vars := value.(string)
				if len(vars) > 0 {
					if len(root.Variables) > 0 {
						return root, fmt.Errorf("%s.inventory: invalid syntax: both variables and variables_json are set", key)
					}

					variables, err := variablesJSONToVariablesMap(vars)
					if err != nil {
						return root, fmt.Errorf("%s.inventory: %s", key, err)
					}
					root.Variables = variables
				}
			}
		}
		return root, nil
	default:
		return root, fmt.Errorf("Unsupported type: %T", v)
	}
}

// ResolvePath checks if a path exists.
func ResolvePath(path string) (string, error) {
	expandedPath, _ := homedir.Expand(path)
	if _, err := os.Stat(expandedPath); err == nil {
		return filepath.Clean(expandedPath), nil
	}
	return "", fmt.Errorf("Ansible module not found at path: [%s]", path)
}

// ResolveDirectory checks if a path exists and is a directory.
func ResolveDirectory(path string) (string, error) {
	expandedPath, _ := homedir.Expand(path)
	if stat, err := os.Stat(expandedPath); err == nil {
		if !stat.IsDir() {
			return "", fmt.Errorf("Path [%s] must be a directory", path)
		}
		return filepath.Clean(expandedPath), nil
	}
	return "", fmt.Errorf("Ansible module not found at path: [%s]", path)
}

func stringToTypeMap(block string, label string) (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(block), &m); err != nil {
		return nil, fmt.Errorf("%s: %s", label, err.Error())
	}
	return m, nil
}

func listOfStringToListOfMap(blocks []string, label string) (output []map[string]interface{}, errs []error) {
	output = make([]map[string]interface{}, 0)
	for i, block := range blocks {
		_map, err := stringToTypeMap(block, fmt.Sprintf("%s[%d]", label, i));
		if err != nil {
			errs = append(errs, err)
		}
		output = append(output, _map)
	}
	return output, errs
}

func variablesJSONToVariablesMap(varsJSON string) (InventoryVariables, error) {
	var m InventoryVariables
	if err := json.Unmarshal([]byte(varsJSON), &m); err != nil {
		return nil, err
	}
	return m, nil
}


func ListOfInventoryHostsToOrderedMap(hosts []InventoryHost) *orderedmap.OrderedMap {
	result := orderedmap.New()
	for _, host := range hosts {
		vars := host.Variables
		if host.AnsibleHost != "" {
			vars["ansible_host"] = host.AnsibleHost
		}
		result.Set(host.Alias, host.Variables)
	}
	return result
}

func ListOfInventoryGroupsToMap(groups []InventoryGroup) InventoryJSONGroupList {
	result := InventoryJSONGroupList{}
	for _, group := range groups {
		groupJSON := InventoryJSONGroup{}

		if len(group.Hosts) > 0 {
			groupJSON.Hosts = ListOfInventoryHostsToOrderedMap(group.Hosts)
		}

		if len(group.Variables) > 0 {
			groupJSON.Vars = group.Variables
		}

		if len(group.Children) > 0 {
			groupJSON.Children = ListOfInventoryGroupsToMap(group.Children)
		}

		result[group.Name] = groupJSON

	}
	return result
}