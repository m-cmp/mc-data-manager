// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "http://cloud-barista.github.io",
            "email": "contact-to-cloud-barista@googlegroups.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/generate/linux": {
            "post": {
                "description": "Generate test data on on-premise Linux.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "[On-premise] Test Data Generation"
                ],
                "summary": "Generate test data on on-premise Linux",
                "parameters": [
                    {
                        "description": "Parameters required to generate test data",
                        "name": "RequestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.GenDataParams"
                        }
                    },
                    {
                        "type": "file",
                        "description": "Parameters required to generate test data",
                        "name": "CredentialGCP",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully generated test data",
                        "schema": {
                            "$ref": "#/definitions/controllers.GenerateLinuxPostHandlerResponseBody"
                        }
                    },
                    "400": {
                        "description": "Invalid Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.GenerateLinuxPostHandlerResponseBody"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.GenerateLinuxPostHandlerResponseBody"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.GenDataParams": {
            "type": "object",
            "properties": {
                "accessKey": {
                    "type": "string"
                },
                "bucket": {
                    "type": "string"
                },
                "checkCSV": {
                    "type": "string"
                },
                "checkGIF": {
                    "type": "string"
                },
                "checkJSON": {
                    "type": "string"
                },
                "checkPNG": {
                    "type": "string"
                },
                "checkSQL": {
                    "type": "string"
                },
                "checkServerJSON": {
                    "type": "string"
                },
                "checkServerSQL": {
                    "type": "string"
                },
                "checkTXT": {
                    "type": "string"
                },
                "checkXML": {
                    "type": "string"
                },
                "checkZIP": {
                    "type": "string"
                },
                "databaseName": {
                    "type": "string"
                },
                "endpoint": {
                    "type": "string"
                },
                "host": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                },
                "port": {
                    "type": "string"
                },
                "projectId": {
                    "type": "string"
                },
                "provider": {
                    "type": "string"
                },
                "region": {
                    "type": "string"
                },
                "secretKey": {
                    "type": "string"
                },
                "sizeCSV": {
                    "type": "string"
                },
                "sizeGIF": {
                    "type": "string"
                },
                "sizeJSON": {
                    "type": "string"
                },
                "sizePNG": {
                    "type": "string"
                },
                "sizeSQL": {
                    "type": "string"
                },
                "sizeServerJSON": {
                    "type": "string"
                },
                "sizeServerSQL": {
                    "type": "string"
                },
                "sizeTXT": {
                    "type": "string"
                },
                "sizeXML": {
                    "type": "string"
                },
                "sizeZIP": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "controllers.GenerateLinuxPostHandlerResponseBody": {
            "type": "object",
            "properties": {
                "Error": {
                    "type": "string"
                },
                "Result": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "latest",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "CM-DataMold REST API",
	Description:      "CM-DataMold REST API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
