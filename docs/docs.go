// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/changes": {
            "get": {
                "description": "Retrieves changes from CouchDB using a specified filter",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "changes"
                ],
                "summary": "Get changes from CouchDB with a filter",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Address",
                        "name": "address",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Age",
                        "name": "age",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Failed to retrieve changes.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/document/{docID}": {
            "put": {
                "description": "Updates an existing document in the CouchDB database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "document"
                ],
                "summary": "Update an existing document",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Document ID",
                        "name": "docID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Document Data",
                        "name": "document",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Document updated successfully",
                        "schema": {
                            "$ref": "#/definitions/main.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/main.Response"
                        }
                    },
                    "404": {
                        "description": "Document not found",
                        "schema": {
                            "$ref": "#/definitions/main.Response"
                        }
                    },
                    "500": {
                        "description": "Failed to update document",
                        "schema": {
                            "$ref": "#/definitions/main.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a document from the CouchDB database",
                "tags": [
                    "document"
                ],
                "summary": "Delete a document",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Document ID",
                        "name": "docID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Document deleted successfully",
                        "schema": {
                            "$ref": "#/definitions/main.DeleteResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Document not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to delete document",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/document/{id}": {
            "get": {
                "description": "Retrieves a specific document from the CouchDB student database by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "document"
                ],
                "summary": "Get a document by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Document ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Document retrieved successfully.",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Document not found.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to retrieve document.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/documents": {
            "get": {
                "description": "Retrieves all documents from the CouchDB student database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "document"
                ],
                "summary": "Get all documents",
                "responses": {
                    "200": {
                        "description": "Documents retrieved successfully.",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "additionalProperties": true
                            }
                        }
                    },
                    "500": {
                        "description": "Failed to retrieve documents.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/file/{docID}/{filename}": {
            "get": {
                "description": "Retrieves an attachment from a CouchDB document",
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "file"
                ],
                "summary": "Get a file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Document ID",
                        "name": "docID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Filename",
                        "name": "filename",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "File downloaded successfully",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "File not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to retrieve file",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/insert": {
            "post": {
                "description": "Inserts a new document into the CouchDB",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "document"
                ],
                "summary": "Insert a document",
                "parameters": [
                    {
                        "description": "Document",
                        "name": "document",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Document inserted successfully.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Failed to decode JSON.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to insert document.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/upload": {
            "post": {
                "description": "Upload a file to CouchDB as an attachment",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Uploads a file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File to upload",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Document ID",
                        "name": "docID",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "File uploaded successfully.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.DeleteResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "main.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "rev": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Student API",
	Description:      "This is a simple API to interact with CouchDB and perform CRUD operations.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
