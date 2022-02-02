package genvaluestring

import (
	"fmt"
	"gogen3/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath     string
	ValueStringName string
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a valueobject with simple string type\n" +
			"   gogen genvaluestring OrderID\n" +
			"     'OrderID' is an valueobject name\n" +
			"\n")
		return err
	}

	domainName := utils.GetDefaultDomain()
	valueStringName := inputs[0]

	obj := &ObjTemplate{
		PackagePath:     utils.GetPackagePath(),
		ValueStringName: valueStringName,
	}

	fileRenamer := map[string]string{
		"valuestringname": utils.SnakeCase(valueStringName),
		"domainname":      utils.LowerCase(domainName),
	}

	err := utils.CreateEverythingExactly("templates/", "valuestring", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	return nil

}
