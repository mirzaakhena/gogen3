package genusecase

import (
	"fmt"
	"gogen3/utils"
)

type GenerateUsecase struct {
}

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	UsecaseName string
}

func (GenerateUsecase) Run(inputs ...string) error {

	if len(inputs) < 2 {
		err := fmt.Errorf("" +
			"# Create a new usecase\n" +
			"\n" +
			"    gogen usecase orderservice CreateOrder\n" +
			"\n" +
			"    'orderservice' is a domain name\n" +
			"    'CreateOrder' is an usecase name\n" +
			"\n")

		return err
	}

	domainName := inputs[0]
	usecaseName := inputs[1]

	obj := &ObjTemplate{
		PackagePath: utils.GetPackagePath(),
		UsecaseName: usecaseName,
	}

	fileRenamer := map[string]string{
		"usecasename": utils.LowerCase(usecaseName),
		"domainname":  utils.LowerCase(domainName),
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