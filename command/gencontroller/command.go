package gencontroller

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"gogen3/utils"
	"os"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath    string
	DomainName     string
	ControllerName string
	UsecaseName    string
	DriverName     string
}

func Run(inputs ...string) error {

	if len(inputs) < 3 {

		err := fmt.Errorf("\n" +
			"   # Create a controller with defined web framework or other handler\n" +
			"   gogen controller orderservice restapi CreateOrder gin\n" +
			"     'orderservice' is an domain name\n" +
			"     'restapi'      is a gateway name\n" +
			"     'CreateOrder'  is an usecase name\n" +
			"     'gin'          is a sample webframewrok. You may try the other one like: nethttp, echo, and gorilla\n" +
			"\n" +
			"   # Create a controller with gin as default web framework\n" +
			"   gogen controller orderservice restapi CreateOrder\n" +
			"     'orderservice' is an domain name\n" +
			"     'restapi'      is a gateway name\n" +
			"     'CreateOrder'  is an usecase name\n" +
			"\n")

		return err
	}

	domainName := inputs[0]
	controllerName := inputs[1]
	usecaseName := inputs[2]

	obj := ObjTemplate{
		PackagePath:    utils.GetPackagePath(),
		DomainName:     domainName,
		ControllerName: controllerName,
		UsecaseName:    usecaseName,
		DriverName:     "gin",
	}

	if len(inputs) >= 4 {
		obj.DriverName = utils.LowerCase(inputs[3])
	}

	fileRenamer := map[string]string{
		"controllername": utils.LowerCase(controllerName),
		"domainname":     utils.LowerCase(domainName),
		"usecasename":    utils.LowerCase(usecaseName),
	}

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	err = utils.CreateEverythingExactly("templates/controllers/", obj.DriverName, fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	//inject inport to struct
	//type Controller struct {
	//  Router            gin.IRouter
	//  CreateOrderInport createorder.Inport <----- here
	//}
	{
		templateCode, err := getRouterInportTemplate(obj.DriverName)
		if err != nil {
			return err
		}

		templateWithData, err := utils.PrintTemplate(string(templateCode), obj)
		if err != nil {
			return err
		}

		dataInBytes, err := injectInportToStruct(obj, templateWithData)
		if err != nil {
			return err
		}

		// reformat router.go
		err = utils.Reformat(obj.getControllerRouterFileName(), dataInBytes)
		if err != nil {
			return err
		}
	}

	// inject router for register
	//func (r *Controller) RegisterRouter() {
	//  r.Router.POST("/createorder", r.authorized(), r.createOrderHandler(r.CreateOrderInport)) <-- here
	//}
	{
		templateCode, err := getRouterRegisterTemplate(obj.DriverName)

		templateWithData, err := utils.PrintTemplate(string(templateCode), obj)
		if err != nil {
			return err
		}

		dataInBytes, err := injectRouterBind(obj, templateWithData)
		if err != nil {
			return err
		}

		// reformat router.go
		err = utils.Reformat(obj.getControllerRouterFileName(), dataInBytes)
		if err != nil {
			return err
		}
	}

	return nil

}

func getRouterInportTemplate(driverName string) ([]byte, error) {
	path := fmt.Sprintf("templates/controllers/%s/domain_${domainname}/controller/${controllername}/~inject-router-inport._go", driverName)
	return utils.AppTemplates.ReadFile(path)
}

func getRouterRegisterTemplate(driverName string) ([]byte, error) {
	path := fmt.Sprintf("templates/controllers/%s/domain_${domainname}/controller/${controllername}/~inject-router-register._go", driverName)
	return utils.AppTemplates.ReadFile(path)
}

func injectInportToStruct(obj ObjTemplate, templateWithData string) ([]byte, error) {

	inportLine, err := getInportLine(obj)
	if err != nil {
		return nil, err
	}

	controllerFile := obj.getControllerRouterFileName()

	file, err := os.Open(controllerFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	line := 0
	for scanner.Scan() {
		row := scanner.Text()

		if line == inportLine-1 {
			buffer.WriteString(templateWithData)
			buffer.WriteString("\n")
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
		line++
	}

	return buffer.Bytes(), nil
}

func getInportLine(obj ObjTemplate) (int, error) {

	controllerFile := obj.getControllerRouterFileName()

	inportLine := 0
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, controllerFile, nil, parser.ParseComments)
	if err != nil {
		return 0, err
	}

	// loop the outport for imports
	for _, decl := range astFile.Decls {

		if gen, ok := decl.(*ast.GenDecl); ok {

			if gen.Tok != token.TYPE {
				continue
			}

			for _, specs := range gen.Specs {

				ts, ok := specs.(*ast.TypeSpec)
				if !ok {
					continue
				}

				if iStruct, ok := ts.Type.(*ast.StructType); ok {

					// check the specific struct name
					if ts.Name.String() != "Controller" {
						continue
					}

					inportLine = fset.Position(iStruct.Fields.Closing).Line
					return inportLine, nil
				}

			}

		}

	}

	return 0, fmt.Errorf(" Controller struct not found")

}

func getBindRouterLine(obj ObjTemplate) (int, error) {

	controllerFile := obj.getControllerRouterFileName()

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, controllerFile, nil, parser.ParseComments)
	if err != nil {
		return 0, err
	}
	routerLine := 0
	for _, decl := range astFile.Decls {

		if gen, ok := decl.(*ast.FuncDecl); ok {

			if gen.Recv == nil {
				continue
			}

			startExp, ok := gen.Recv.List[0].Type.(*ast.StarExpr)
			if !ok {
				continue
			}

			if startExp.X.(*ast.Ident).String() != "Controller" {
				continue
			}

			if gen.Name.String() != "RegisterRouter" {
				continue
			}

			routerLine = fset.Position(gen.Body.Rbrace).Line
			return routerLine, nil
		}

	}
	return 0, fmt.Errorf("register router Not found")
}

func injectRouterBind(obj ObjTemplate, templateWithData string) ([]byte, error) {

	controllerFile := obj.getControllerRouterFileName()

	routerLine, err := getBindRouterLine(obj)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(controllerFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	line := 0
	for scanner.Scan() {
		row := scanner.Text()

		if line == routerLine-1 {
			buffer.WriteString(templateWithData)
			buffer.WriteString("\n")
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
		line++
	}

	return buffer.Bytes(), nil

}

// getControllerRouterFileName ...
func (o ObjTemplate) getControllerRouterFileName() string {
	return fmt.Sprintf("domain_%s/controller/%s/router.go", utils.LowerCase(o.DomainName), utils.LowerCase(o.ControllerName))
}
