package genopenapi

import (
	"encoding/json"
	"fmt"
	"gogen3/utils"
	"io/ioutil"
)

// ObjTemplate ...
type ObjTemplate struct {
	GomodPath     string
	DefaultDomain string
}

type OASSchema struct {
	Type       string `json:"type"`
	Properties []any  `json:"properties"`
	Required   bool   `json:"required"`
}

type OASParameter struct {
	Name        string    `json:"name"`
	In          string    `json:"in"`
	Description string    `json:"description"`
	Required    bool      `json:"required"`
	Schema      OASSchema `json:"schema"`
}

type OASExample struct {
	Description   string `json:"description"`
	Summary       string `json:"summary"`
	Value         any    `json:"value"`
	ExternalValue string `json:"externalValue"`
}

type OASMediaType struct {
	Schema   string                `json:"schema,omitempty"`
	Example  any                   `json:"example,omitempty"`
	Examples map[string]OASExample `json:"examples,omitempty"`
}

type OASRequestBody struct {
	Description string         `json:"description,omitempty"`
	Content     map[string]any `json:"content"`
	Required    bool           `json:"required,omitempty"`
}

type OASResponse struct {
	Description string                  `json:"description"`
	Content     map[string]OASMediaType `json:"content"`
}

type OASOperationObject struct {
	Tags        []string               `json:"tags"`
	Summary     string                 `json:"summary"`
	OperationID string                 `json:"operationId"`
	Parameters  []OASParameter         `json:"parameters"`
	RequestBody OASRequestBody         `json:"requestBody"`
	Responses   map[string]OASResponse `json:"responses"`
}

type OASPathItem struct {
	Get         OASOperationObject `json:"get,omitempty"`
	Post        OASOperationObject `json:"post,omitempty"`
	Put         OASOperationObject `json:"put,omitempty"`
	Delete      OASOperationObject `json:"delete,omitempty"`
	Summary     string             `json:"summary,omitempty"`
	Description string             `json:"description,omitempty"`
	Ref         string             `json:"$ref"`
	Servers     OASServer          `json:"servers,omitempty"`
	Parameters  OASParameter       `json:"parameters,omitempty"`
}

type OASExternalDocumentation struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

type OASServerVariable struct {
	Description string   `json:"description"`
	Default     string   `json:"default"`
	Enum        []string `json:"enum,omitempty"`
}

//type OpenAPIVariables struct {
//	Protocol    OASServerVariable `json:"protocols"`
//	Environment OASServerVariable `json:"environment"`
//	Port        OASServerVariable `json:"port"`
//	BasePath    OASServerVariable `json:"basePath"`
//}

type OASServer struct {
	Description string                       `json:"description"`
	Url         string                       `json:"url"`
	Variables   map[string]OASServerVariable `json:"variables"`
}

type OpenAPILicense struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type OpenAPIContact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Url   string `json:"url"`
}

type OpenAPIInfo struct {
	Title          string         `json:"title"`
	Version        string         `json:"version"`
	Description    string         `json:"description"`
	TermsOfService string         `json:"termsOfService"`
	Contact        OpenAPIContact `json:"contact"`
	License        OpenAPILicense `json:"license"`
}

type OpenAPIRoot struct {
	OpenAPI      string                   `json:"openapi"`
	Info         OpenAPIInfo              `json:"info"`
	Servers      []OASServer              `json:"servers"`
	ExternalDocs OASExternalDocumentation `json:"externalDocs"`
	Paths        map[string]OASPathItem
}

func Run(inputs ...string) error {

	domainName := utils.GetDefaultDomain()

	data := OpenAPIRoot{
		OpenAPI: "3.0.3",
		Info: OpenAPIInfo{
			Title:          "Application Restful API",
			Version:        "1",
			Description:    "Application Restful API",
			TermsOfService: "https://gogen.com/tnc",
			Contact: OpenAPIContact{
				Name:  "Mirza Akhena",
				Email: "mirza.akhena@gmail.com",
				Url:   "mirzaakhena.com",
			},
			License: OpenAPILicense{
				Name: "APACHE 2.0",
				Url:  "https://www.apache.org/licenses/LICENSE-2.0",
			},
		},
		Servers: []OASServer{
			{
				Description: fmt.Sprintf("%s", domainName),
				Url:         "https://localhost:8080/api/v1",
				Variables: map[string]OASServerVariable{
					"protocol": {
						Description: "Protocol",
						Default:     "http",
						Enum:        []string{"http", "https"},
					},
					"environment": {
						Description: "Environment",
						Default:     "localhost",
						Enum:        []string{"localhost", "dev", "qa", "prod"},
					},
					"port": {
						Description: "Port",
						Default:     "8080",
						Enum:        []string{"8080", "8081", "8082", "80", "443"},
					},
					"basepath": {
						Description: "Base Path",
						Default:     "v1",
						Enum:        nil,
					},
				},
				//Variables: OpenAPIVariables{
				//	Environment: OASServerVariable{
				//		Description: "Environment",
				//		Default:     "localhost",
				//		Enum:        []string{"localhost", "dev", "qa", "prod"},
				//	},

				//},
			},
		},
		ExternalDocs: OASExternalDocumentation{
			Description: fmt.Sprintf("Documentation for %s", domainName),
			URL:         fmt.Sprintf("https://bitbucket.org/%s", domainName),
		},
		Paths: map[string]OASPathItem{},
	}

	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("domain_%s/openapi.json", domainName), file, 0644)
	if err != nil {
		return err
	}

	return nil
}
