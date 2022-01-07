package gengateway

import (
	"fmt"
	"gogen3/utils"
	"io/ioutil"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	DomainName  string
	GatewayName string
	UsecaseName *string
	Methods     utils.OutportMethods
}

func Run(inputs ...string) error {

	if len(inputs) < 2 {
		err := fmt.Errorf("\n" +
			"   # Create a gateway for specific usecase\n" +
			"   gogen gateway orderservice inmemory CreateOrder\n" +
			"     'orderservice' is an domain name\n" +
			"     'inmemory'     is a gateway name\n" +
			"     'CreateOrder'  is an usecase name\n" +
			"\n" +
			"   # Create a gateway for all usecases\n" +
			"   gogen gateway orderservice inmemory\n" +
			"     'orderservice' is an domain name\n" +
			"     'inmemory'     is a gateway name\n" +
			"\n")

		return err
	}

	domainName := inputs[0]
	gatewayName := inputs[1]

	obj := ObjTemplate{
		PackagePath: utils.GetPackagePath(),
		GatewayName: gatewayName,
		DomainName:  utils.LowerCase(domainName),
		UsecaseName: nil,
	}

	if len(inputs) >= 3 {
		obj.UsecaseName = &inputs[2]
	}

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	var notExistingMethod utils.OutportMethods

	if obj.UsecaseName == nil {

		var folders []string
		fileInfo, err := ioutil.ReadDir(fmt.Sprintf("domain_%s/usecase", domainName))
		if err != nil {
			return err
		}

		uniqueMethodMap := map[string]int{}

		for _, file := range fileInfo {

			folders = append(folders, file.Name())

			em, err := createGatewayImpl(file.Name(), obj)
			if err != nil {
				return err
			}

			for _, method := range em {

				if _, exist := uniqueMethodMap[method.MethodName]; exist {
					continue
				}

				notExistingMethod = append(notExistingMethod, method)

				uniqueMethodMap[method.MethodName] = 1
			}
		}

	} else {

		em, err := createGatewayImpl(*obj.UsecaseName, obj)
		if err != nil {
			return err
		}

		for _, method := range em {
			notExistingMethod = append(notExistingMethod, method)
		}

	}

	gatewayCode, err := getGatewayMethodTemplate()
	if err != nil {
		return err
	}

	// we will only inject the non existing method
	obj.Methods = notExistingMethod

	templateHasBeenInjected, err := utils.PrintTemplate(string(gatewayCode), obj)
	if err != nil {
		return err
	}

	gatewayFilename := fmt.Sprintf("domain_%s/gateway/%s/gateway.go", domainName, gatewayName)

	bytes, err := injectToGateway(gatewayFilename, templateHasBeenInjected)
	if err != nil {
		return err
	}

	// reformat outport.go
	err = utils.Reformat(gatewayFilename, bytes)
	if err != nil {
		return err
	}

	return nil

}

func createGatewayImpl(usecaseName string, obj ObjTemplate) (utils.OutportMethods, error) {
	outportMethods, err := utils.NewOutportMethods(obj.DomainName, usecaseName)
	if err != nil {
		return nil, err
	}

	obj.Methods = outportMethods
	err = utils.CreateEverythingExactly("templates/", "gateway", map[string]string{
		"gatewayname": utils.LowerCase(obj.GatewayName),
		"domainname":  obj.DomainName,
	}, obj, utils.AppTemplates)
	if err != nil {
		return nil, err
	}

	gatewayRootFolderName := fmt.Sprintf("domain_%s/gateway/%s", obj.DomainName, utils.LowerCase(obj.GatewayName))

	// file gateway impl file is already exist, we want to inject non existing method
	existingFunc, err := utils.NewOutportMethodImpl("gateway", gatewayRootFolderName)
	if err != nil {
		return nil, err
	}

	// collect the only methods that has not added yet
	notExistingMethod := utils.OutportMethods{}
	for _, m := range outportMethods {
		if _, exist := existingFunc[m.MethodName]; !exist {
			notExistingMethod = append(notExistingMethod, m)
		}
	}
	return notExistingMethod, nil
}

// getGatewayMethodTemplate ...
func getGatewayMethodTemplate() ([]byte, error) {
	return utils.AppTemplates.ReadFile("templates/gateway/domain_${domainname}/gateway/${gatewayname}/~inject._go")
}

func injectToGateway(gatewayFilename, injectedCode string) ([]byte, error) {
	return utils.InjectCodeAtTheEndOfFile(gatewayFilename, injectedCode)
}
