package geninit

import (
	"bufio"
	"fmt"
	"gogen3/utils"
	"os"
	"strings"
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
.idea/
.DS_Store
config.json
*.app
*.exe
*.log
*.db
*/node_modules/
	`
	_, err = utils.WriteFileIfNotExist(gitignoreContent, "./.gitignore", struct{}{})
	if err != nil {
		return err
	}

	inFile, err := os.Open(".gogen/domain")
	if err != nil {
		return err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		domainName := strings.TrimSpace(scanner.Text())
		if domainName == "" {
			continue
		}
		if strings.HasPrefix(domainName, "-") {
			domainName = strings.ReplaceAll(domainName, "-", "")
		}
		domainName = strings.ToLower(domainName)
		_, err := utils.CreateFolderIfNotExist(fmt.Sprintf("domain_%s", domainName))
		if err != nil {
			return err
		}

	}

	fmt.Printf("open .gogen/domain file and add all your domains name. Then run 'gogen init' again to create domain folders\n")

	return nil

}
