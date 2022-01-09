package gencrud

import (
	"fmt"
	"go/token"
	"gogen3/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	EntityName  string
	DomainName  string
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
		DomainName:  domainName,
	}

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(domainName),
		"entityname": utils.LowerCase(entityName),
	}

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	err = utils.CreateEverythingExactly("templates/", "crud", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	// inject to main.go
	{
		fset := token.NewFileSet()
		utils.InjectToMain(fset, fmt.Sprintf("App%s", entityName))
	}

	return nil

}
