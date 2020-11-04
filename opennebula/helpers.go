package opennebula

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/OpenNebula/one/src/oca/go/src/goca/errors"
	"github.com/OpenNebula/one/src/oca/go/src/goca/schemas/shared"

	"github.com/hashicorp/terraform/helper/schema"
)

func inArray(val string, array []string) (index int) {
	var ok bool
	for i := range array {
		if ok = array[i] == val; ok {
			return i
		}
	}
	return -1
}

// appendTemplate add attribute and value to an existing string
func appendTemplate(template, attribute, value string) string {
	return fmt.Sprintf("%s\n%s = \"%s\"", template, attribute, value)
}

func ArrayToString(list []interface{}, delim string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(list)), delim), "[]")
}

func StringToLockLevel(str string, lock *shared.LockLevel) error {
	if str == "USE" {
		*lock = shared.LockUse
		return nil
	}
	if str == "MANAGE" {
		*lock = shared.LockManage
		return nil
	}
	if str == "ADMIN" {
		*lock = shared.LockAdmin
		return nil
	}
	if str == "ALL" {
		*lock = shared.LockAll
		return nil
	}
	return fmt.Errorf("Unexpected Lock level %s", str)
}

func LockLevelToString(lock int) string {
	if lock == 1 {
		return "USE"
	}
	if lock == 2 {
		return "MANAGE"
	}
	if lock == 3 {
		return "ADMIN"
	}
	if lock == 4 {
		return "ALL"
	}
	return ""
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// NoExists indicate if an entity exists in checking the error code returned from an Info call
func NoExists(err error) bool {

	respErr, ok := err.(*errors.ResponseError)

	// expected case, the entity does not exists so we doesn't return an error
	if ok && respErr.Code == errors.OneNoExistsError {
		return true
	}

	return false
}

// IDSet is a set implementation that allow to manipulate schema configs based on IDs
type IDSet struct{ *schema.Set }

func NewIDSet(IDs ...interface{}) *IDSet {
	return &IDSet{
		schema.NewSet(schema.HashInt, IDs),
	}
}

// InsertConfigIDs insert ID from config via it's name attrName.
// attrName should be the name of an attribute of type int
func (s *IDSet) InsertConfigIDs(schemaList []interface{}, attrName string) {
	for _, item := range schemaList {

		mapItem := item.(map[string]interface{})

		id := mapItem[attrName].(int)

		if id < 0 {
			continue
		}

		s.Add(id)
	}
}

// DiffConfigIDs achieve a partial diff based on ID (via attrName) and return slice of config that only appear on ref side
// attrName should be an ID of type int
func (s *IDSet) DiffConfigIDs(schemaList []interface{}, attrName string) []interface{} {

	partialDiff := make([]interface{}, 0)

	for _, item := range schemaList {

		mapItem := item.(map[string]interface{})

		id := mapItem[attrName].(int)

		if id < 0 {
			continue
		}

		if !s.Contains(id) {
			partialDiff = append(partialDiff, mapItem)
		}
	}

	return partialDiff
}

// partial diff that return a slice of config (map[string]interface{}) that only appear on ref side based on the attrName attribute.
// attrName should be an ID of type int
func diffIDsConfig(refVecs, vecs []interface{}, attrName string) []interface{} {

	set := NewIDSet()

	set.InsertConfigIDs(vecs, attrName)

	return set.DiffConfigIDs(refVecs, attrName)
}
