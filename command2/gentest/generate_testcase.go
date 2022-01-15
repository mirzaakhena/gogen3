package gentest

import (
	"fmt"
	"gogen3/utils"
)

type GenerateTest struct {
}

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	DomainName  string
	UsecaseName string
	TestName    string
	Methods     utils.OutportMethods
}

func (GenerateTest) Run(inputs ...string) error {

	if len(inputs) < 3 {
		err := fmt.Errorf("\n" +
			"   # Create a test case file for current usecase\n" +
			"   gogen test orderservice normal CreateOrder\n" +
			"     'orderservice' is a domain name\n" +
			"     'normal'       is a test case name\n" +
			"     'CreateOrder'  is an usecase name\n" +
			"\n")

		return err
	}

	domainName := inputs[0]
	testName := inputs[1]
	usecaseName := inputs[2]

	obj := ObjTemplate{
		PackagePath: utils.GetPackagePath(),
		UsecaseName: usecaseName,
		TestName:    testName,
		DomainName:  domainName,
		Methods:     nil,
	}

	outportMethods, err := utils.NewOutportMethods(domainName, usecaseName)
	if err != nil {
		return err
	}

	obj.Methods = outportMethods

	fileRenamer := map[string]string{
		"domainname":  utils.LowerCase(domainName),
		"usecasename": utils.LowerCase(usecaseName),
		"testname":    utils.LowerCase(testName),
	}
	err = utils.CreateEverythingExactly("templates/", "test", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	// file gateway impl file is already exist, we want to inject non existing method
	structName := fmt.Sprintf("mockOutport%s", utils.PascalCase(testName))
	rootFolderName := fmt.Sprintf("domain_%s/usecase/%s", utils.LowerCase(domainName), utils.LowerCase(usecaseName))
	existingFunc, err := utils.NewOutportMethodImpl(structName, rootFolderName)
	if err != nil {
		return err
	}

	// collect the only methods that has not added yet
	notExistingMethod := utils.OutportMethods{}
	for _, m := range outportMethods {
		if _, exist := existingFunc[m.MethodName]; !exist {
			notExistingMethod = append(notExistingMethod, m)
		}
	}

	gatewayCode, err := getTestMethodTemplate()
	if err != nil {
		return err
	}

	obj.Methods = notExistingMethod

	templateHasBeenInjected, err := utils.PrintTemplate(string(gatewayCode), obj)
	if err != nil {
		return err
	}

	bytes, err := obj.injectToTest(templateHasBeenInjected)
	if err != nil {
		return err
	}

	// reformat outport._go
	err = utils.Reformat(obj.getTestFileName(), bytes)
	if err != nil {
		return err
	}

	return nil

}

// getTestMethodTemplate ...
func getTestMethodTemplate() ([]byte, error) {
	testfileName := fmt.Sprintf("templates/test/domain_${domainname}/usecase/${usecasename}/~inject._go")
	return utils.AppTemplates.ReadFile(testfileName)
}

func (o ObjTemplate) injectToTest(injectedCode string) ([]byte, error) {
	return utils.InjectCodeAtTheEndOfFile(o.getTestFileName(), injectedCode)
}

// getTestFileName ...
func (o ObjTemplate) getTestFileName() string {
	return fmt.Sprintf("domain_%s/usecase/%s/testcase_%s_test.go", utils.LowerCase(o.DomainName), utils.LowerCase(o.UsecaseName), utils.LowerCase(o.TestName))
}
