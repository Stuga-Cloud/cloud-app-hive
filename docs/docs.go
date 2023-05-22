// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/applications": {
            "post": {
                "description": "creates in database and deploys an application on the cloud",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Applications"
                ],
                "summary": "Creates in database and deploys an application",
                "operationId": "create-and-deploy-application",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Create Application Request",
                        "name": "createApplicationRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.CreateApplicationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.CreateApplicationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.ApiError"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "checks the health of the API",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health"
                ],
                "summary": "Health check endpoint",
                "operationId": "health-check",
                "responses": {
                    "200": {
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.ApplicationContainerSpecifications": {
            "type": "object",
            "properties": {
                "cpu_limit": {
                    "$ref": "#/definitions/domain.ContainerLimit"
                },
                "memory_limit": {
                    "$ref": "#/definitions/domain.ContainerLimit"
                },
                "storage_limit": {
                    "$ref": "#/definitions/domain.ContainerLimit"
                }
            }
        },
        "domain.ApplicationEnvironmentVariable": {
            "type": "object",
            "required": [
                "name",
                "value"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "domain.ApplicationScalabilitySpecifications": {
            "type": "object",
            "required": [
                "maximum_instance_count",
                "minimum_instance_count",
                "replicas"
            ],
            "properties": {
                "is_auto_scaled": {
                    "type": "boolean"
                },
                "maximum_instance_count": {
                    "type": "integer"
                },
                "minimum_instance_count": {
                    "type": "integer"
                },
                "replicas": {
                    "type": "integer"
                }
            }
        },
        "domain.ApplicationSecret": {
            "type": "object",
            "required": [
                "name",
                "value"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "domain.ApplicationType": {
            "type": "string",
            "enum": [
                "SINGLE_INSTANCE",
                "LOAD_BALANCED"
            ],
            "x-enum-varnames": [
                "SingleInstance",
                "LoadBalanced"
            ]
        },
        "domain.ContainerLimit": {
            "type": "object",
            "properties": {
                "unit": {
                    "enum": [
                        "KB",
                        "MB",
                        "GB",
                        "TB"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/domain.LimitUnit"
                        }
                    ]
                },
                "value": {
                    "type": "integer"
                }
            }
        },
        "domain.LimitUnit": {
            "type": "string",
            "enum": [
                "KB",
                "MB",
                "GB",
                "TB"
            ],
            "x-enum-varnames": [
                "KB",
                "MB",
                "GB",
                "TB"
            ]
        },
        "errors.ApiError": {
            "type": "object",
            "properties": {
                "context": {},
                "date": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "requests.CreateApplicationRequest": {
            "type": "object",
            "required": [
                "image",
                "name",
                "namespace_id",
                "port",
                "user_id"
            ],
            "properties": {
                "application_type": {
                    "enum": [
                        "SINGLE_INSTANCE",
                        "LOAD_BALANCED"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/domain.ApplicationType"
                        }
                    ]
                },
                "container_specifications": {
                    "$ref": "#/definitions/domain.ApplicationContainerSpecifications"
                },
                "environment_variables": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.ApplicationEnvironmentVariable"
                    }
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 3
                },
                "namespace_id": {
                    "type": "string"
                },
                "port": {
                    "type": "integer",
                    "maximum": 65535,
                    "minimum": 1
                },
                "scalability_specifications": {
                    "$ref": "#/definitions/domain.ApplicationScalabilitySpecifications"
                },
                "secrets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.ApplicationSecret"
                    }
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "responses.ApplicationResponse": {
            "type": "object",
            "properties": {
                "application_type": {
                    "$ref": "#/definitions/domain.ApplicationType"
                },
                "container_specifications": {
                    "$ref": "#/definitions/domain.ApplicationContainerSpecifications"
                },
                "environment_variables": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.ApplicationEnvironmentVariable"
                    }
                },
                "id": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "namespace_id": {
                    "type": "string"
                },
                "port": {
                    "type": "integer"
                },
                "scalability_specifications": {
                    "$ref": "#/definitions/domain.ApplicationScalabilitySpecifications"
                },
                "secrets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.ApplicationSecret"
                    }
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "responses.CreateApplicationResponse": {
            "type": "object",
            "properties": {
                "application": {
                    "$ref": "#/definitions/responses.ApplicationResponse"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
