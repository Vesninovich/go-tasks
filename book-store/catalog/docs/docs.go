// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Dimas",
            "url": "https://github.com/Vesninovich",
            "email": "dmitry@vesnin.work"
        },
        "license": {
            "name": "ISC",
            "url": "https://www.isc.org/licenses/"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/book": {
            "get": {
                "description": "get books according to query",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "get books",
                "parameters": [
                    {
                        "type": "string",
                        "description": "results start",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "results count",
                        "name": "count",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "book id",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "author id",
                        "name": "author",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "category ids",
                        "name": "categories",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "results",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github.com_Vesninovich_go-tasks_book-store_catalog_rest.apiModel"
                            }
                        }
                    },
                    "400": {
                        "description": "malformed query",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "create book",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Book"
                ],
                "summary": "create book",
                "parameters": [
                    {
                        "description": "book data",
                        "name": "order",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github.com_Vesninovich_go-tasks_book-store_catalog_rest.createAPIModel"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "created book",
                        "schema": {
                            "$ref": "#/definitions/book.Book"
                        }
                    },
                    "400": {
                        "description": "malformed data",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "nested author or category not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "book.Author": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "book.Book": {
            "type": "object",
            "properties": {
                "author": {
                    "$ref": "#/definitions/book.Author"
                },
                "categories": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/book.Category"
                    }
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "book.Category": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "parentID": {
                    "type": "string"
                }
            }
        },
        "github.com_Vesninovich_go-tasks_book-store_catalog_rest.apiModel": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "categories": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "github.com_Vesninovich_go-tasks_book-store_catalog_rest.createAPIModel": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "object",
                    "properties": {
                        "id": {
                            "type": "string"
                        },
                        "name": {
                            "type": "string"
                        }
                    }
                },
                "categories": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "id": {
                                "type": "string"
                            },
                            "name": {
                                "type": "string"
                            },
                            "parentID": {
                                "type": "string"
                            }
                        }
                    }
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "rest.apiModel": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "categories": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "rest.createAPIModel": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "object",
                    "properties": {
                        "id": {
                            "type": "string"
                        },
                        "name": {
                            "type": "string"
                        }
                    }
                },
                "categories": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "id": {
                                "type": "string"
                            },
                            "name": {
                                "type": "string"
                            },
                            "parentID": {
                                "type": "string"
                            }
                        }
                    }
                },
                "name": {
                    "type": "string"
                }
            }
        }
    },
    "tags": [
        {
            "description": "Quering and creating books",
            "name": "Book"
        }
    ]
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.0",
	Host:        "localhost:8002",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Book Store Catalog Service",
	Description: "Service for creating and quering books catalog",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
