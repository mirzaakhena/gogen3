package genvalueobject

import (
	"fmt"
	"gogen3/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath     string
	ValueObjectName string
	FieldNames      []string
}

func Run(inputs ...string) error {

	if len(inputs) < 2 {
		err := fmt.Errorf("\n" +
			"   # Create a valueobject with struct type\n" +
			"   gogen valueobject orderservice FullName FirstName LastName\n" +
			"     'orderservice'             is a domain name\n" +
			"     'FullName'                 is a Value Object name to created\n" +
			"     'FirstName' and 'LastName' is a Fields on Value Object\n" +
			"\n")
		return err
	}

	domainName := inputs[0]
	valueObjectName := inputs[1]

	obj := &ObjTemplate{
		PackagePath:     utils.GetPackagePath(),
		ValueObjectName: valueObjectName,
		FieldNames:      inputs[2:],
	}

	fileRenamer := map[string]string{
		"valueobjectname": utils.SnakeCase(valueObjectName),
		"domainname":      utils.LowerCase(domainName),
	}

	err := utils.CreateEverythingExactly("templates/", "valueobject", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	return nil

}
