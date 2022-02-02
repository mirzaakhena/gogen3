package genentity

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

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create an entity\n" +
			"   gogen entity Product\n" +
			"     'Product' is an entity name\n" +
			"\n")
		return err
	}

	domainName := utils.GetDefaultDomain()
	entityName := inputs[0]

	obj := &ObjTemplate{
		PackagePath: utils.GetPackagePath(),
		EntityName:  entityName,
	}

	fileRenamer := map[string]string{
		"entityname": utils.SnakeCase(entityName),
		"domainname": utils.LowerCase(domainName),
	}

	err := utils.CreateEverythingExactly("templates/", "entity", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	return nil

}
