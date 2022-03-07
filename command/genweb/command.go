package genweb

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"gogen3/utils"
	"io/ioutil"
	"os"
	"strings"
)

// ObjTemplate ...
type ObjTemplate struct {
	DomainName  string
	UsecaseName string
}

func Run(inputs ...string) error {

	//if len(inputs) < 1 {
	//	err := fmt.Errorf("\n" +
	//		"   # Create a web app\n" +
	//		"   gogen web\n" +
	//		"     'Product' is an existing entity name\n" +
	//		"\n")
	//
	//	return err
	//}

	domainName := utils.GetDefaultDomain()

	//entityName := inputs[0]

	controllerName := "restapi"

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(domainName),
	}

	err := utils.CreateEverythingExactly("templates/web/", "shared", fileRenamer, struct{}{}, utils.AppTemplates)
	if err != nil {
		return err
	}

	controllerFolderName := fmt.Sprintf("domain_%s/controller/%s", domainName, controllerName)

	fileInfo, err := ioutil.ReadDir(controllerFolderName)
	if err != nil {
		return err
	}

	for _, file := range fileInfo {

		if !strings.HasPrefix(file.Name(), "handler_") {
			continue
		}

		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, fmt.Sprintf("%s/%s", controllerFolderName, file.Name()), nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(1)
		}

		//ast.Print(fset, astFile)

		for _, decl := range astFile.Decls {

			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			methodName := utils.PascalCase(funcDecl.Name.String())
			if !strings.HasSuffix(methodName, "Handler") {
				continue
			}

			i := strings.LastIndex(methodName, "Handler")

			usecaseName := methodName[:i]

			fileRenamer := map[string]string{
				"domainname":  utils.LowerCase(domainName),
				"usecasename": utils.LowerCase(usecaseName),
			}

			obj := &ObjTemplate{
				DomainName:  domainName,
				UsecaseName: usecaseName,
			}

			if strings.HasPrefix(utils.LowerCase(usecaseName), "getall") {
				err := utils.CreateEverythingExactly("templates/web/", "getall", fileRenamer, obj, utils.AppTemplates)
				if err != nil {
					return err
				}
			} else if strings.HasPrefix(utils.LowerCase(usecaseName), "get") {
				err := utils.CreateEverythingExactly("templates/web/", "get", fileRenamer, obj, utils.AppTemplates)
				if err != nil {
					return err
				}
			} else {
				err := utils.CreateEverythingExactly("templates/web/", "run", fileRenamer, obj, utils.AppTemplates)
				if err != nil {
					return err
				}
			}

		}

		//break

	}

	return nil

}
