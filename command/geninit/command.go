package geninit

import (
	"fmt"
	"gogen3/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	GomodPath     string
	DefaultDomain string
}

func Run(inputs ...string) error {

	if len(inputs) < 0 {
		err := fmt.Errorf("\n" +
			"   # Initiate gogen project with default input. You may change later under .gogen folder\n" +
			"   gogen init \n" +
			"\n")

		return err
	}

	gomodPath := "your/path/project"
	defaultDomain := "-yourdefaultdomain"

	obj := &ObjTemplate{
		GomodPath:     gomodPath,
		DefaultDomain: defaultDomain,
	}

	fileRenamer := map[string]string{}

	err := utils.CreateEverythingExactly("templates/", "init", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	_, err = utils.CreateFolderIfNotExist(".gogen")
	if err != nil {
		return err
	}

	_, err = utils.WriteFileIfNotExist(defaultDomain, "./.gogen/domain", struct{}{})
	if err != nil {
		return err
	}

	gitignoreContent := `
./idea
*.DS_Store
config.json
*/node_modules/
	`

	_, err = utils.WriteFileIfNotExist(gitignoreContent, "./.gitignore", struct{}{})
	if err != nil {
		return err
	}

	return nil

}
