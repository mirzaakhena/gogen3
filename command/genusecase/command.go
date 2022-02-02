package genusecase

import (
	"fmt"
	"gogen3/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	UsecaseName string
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a new usecase\n" +
			"   gogen usecase CreateOrder\n" +
			"     'CreateOrder' is an usecase name\n" +
			"\n")

		return err
	}

	domainName := utils.GetDefaultDomain()

	usecaseName := inputs[0]

	obj := &ObjTemplate{
		PackagePath: utils.GetPackagePath(),
		UsecaseName: usecaseName,
	}

	fileRenamer := map[string]string{
		"usecasename": utils.LowerCase(usecaseName),
		"domainname":  domainName,
	}

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	err = utils.CreateEverythingExactly("templates/", "usecase", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	return nil

}
