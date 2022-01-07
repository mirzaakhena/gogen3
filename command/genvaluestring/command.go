package genvaluestring

import (
	"fmt"
	"gogen3/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath     string
	ValueStringName string
}

func Run(inputs ...string) error {

	if len(inputs) < 2 {
		err := fmt.Errorf("\n" +
			"   # Create a valueobject with simple string type\n" +
			"   gogen genvaluestring orderservice OrderID\n" +
			"     'orderservice' is a domain name\n" +
			"     'OrderID'      is an valueobject name\n" +
			"\n")
		return err
	}

	domainName := inputs[0]
	valueStringName := inputs[1]

	obj := &ObjTemplate{
		PackagePath:     utils.GetPackagePath(),
		ValueStringName: valueStringName,
	}

	fileRenamer := map[string]string{
		"valuestringname": utils.SnakeCase(valueStringName),
		"domainname":      utils.LowerCase(domainName),
	}

	err := utils.CreateEverythingExactly("templates/", "valuestring", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	//// code your usecase definition here ...
	//_, err := utils.CreateFolderIfNotExist(fmt.Sprintf("domain/%s/model/vo", domainName))
	//if err != nil {
	//	return err
	//}
	//
	//{
	//	outputFile := fmt.Sprintf("domain/%s/model/vo/%s.go", domainName, utils.SnakeCase(obj.ValueStringName))
	//	tem, err := getValueStringTemplate()
	//	if err != nil {
	//		return err
	//	}
	//
	//	_, err = utils.WriteFileIfNotExist(string(tem), outputFile, obj)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil

}

//// getValueStringTemplate ...
//func getValueStringTemplate() ([]byte, error) {
//	return utils.AppTemplates.ReadFile("templates/domain_${domainname}/model/vo/${valuestringname}._go")
//}
