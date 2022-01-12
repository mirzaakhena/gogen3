package generror

import (
	"fmt"
	"go/token"
	"gogen3/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	ErrorName   string
}

func Run(inputs ...string) error {

	if len(inputs) < 2 {
		err := fmt.Errorf("\n" +
			"   # Create an error enum\n" +
			"   gogen error orderservice SomethingGoesWrongError\n" +
			"     'orderservice' is a domain name\n" +
			"     'SomethingGoesWrongError' is an error constant name\n" +
			"\n")

		return err
	}

	domainName := inputs[0]
	errorName := inputs[1]

	obj := ObjTemplate{
		PackagePath: utils.GetPackagePath(),
		ErrorName:   errorName,
	}

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(domainName),
	}

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	err = utils.CreateEverythingExactly("templates/", "errorenum", fileRenamer, struct{}{}, utils.AppTemplates)
	if err != nil {
		return err
	}

	errEnumFile := fmt.Sprintf("domain_%s/model/errorenum/error_enum.go", domainName)

	// inject to error_enum.go
	{
		fset := token.NewFileSet()
		utils.InjectToErrorEnum(fset, errEnumFile, errorName, "ER")
	}

	// reformat outport._go
	err = utils.Reformat(errEnumFile, nil)
	if err != nil {
		return err
	}

	return nil

}
