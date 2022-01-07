package generror

import (
	"fmt"
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

	// TODO nambahin ErrorType dari shared/infrastructure

	err := utils.CreateEverythingExactly("templates/", "errorenum", fileRenamer, struct{}{}, utils.AppTemplates)
	if err != nil {
		return err
	}

	errorLine, err := GetErrorLineTemplate()
	if err != nil {
		return err
	}

	templateHasBeenInjected, err := utils.PrintTemplate(string(errorLine), obj)
	if err != nil {
		return err
	}

	errEnumFile := fmt.Sprintf("domain_%s/model/errorenum/error_enum.go", domainName)

	bytes, err := injectCode(errEnumFile, templateHasBeenInjected)
	if err != nil {
		return err
	}

	// reformat outport._go
	err = utils.Reformat(errEnumFile, bytes)
	if err != nil {
		return err
	}

	return nil

}

// GetErrorLineTemplate ...
func GetErrorLineTemplate() ([]byte, error) {
	return utils.AppTemplates.ReadFile("templates/errorenum/domain_${domainname}/model/errorenum/~inject._go")
}

// InjectCode ...
func injectCode(errEnumFile, templateCode string) ([]byte, error) {
	return utils.InjectCodeAtTheEndOfFile(errEnumFile, templateCode)
}
