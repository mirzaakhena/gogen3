package gencrud

import (
	"fmt"
	"gogen3/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	EntityName  string
}

func Run(inputs ...string) error {

	if len(inputs) < 2 {
		err := fmt.Errorf("\n" +
			"   # Create a new usecase\n" +
			"   gogen crud orderservice Product\n" +
			"     'orderservice' is a domain name\n" +
			"     'Product' is an existing entity name\n" +
			"\n")

		return err
	}

	domainName := inputs[0]

	entityName := inputs[1]

	obj := &ObjTemplate{
		PackagePath: utils.GetPackagePath(),
		EntityName:  entityName,
	}

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(domainName),
		"entityname": utils.LowerCase(entityName),
	}

	err := utils.CreateEverythingExactly("templates/", "usecase", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	return nil

}
