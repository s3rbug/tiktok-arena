// Code generated by swaggo/swag. DO NOT EDIT
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
        "/auth/login": {
            "post": {
                "description": "Login user with given credentials",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "Data to login user",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AuthInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login success",
                        "schema": {
                            "$ref": "#/definitions/models.UserAuthDetails"
                        }
                    },
                    "400": {
                        "description": "Error logging in",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Register new user with given credentials",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "Data to register user",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AuthInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Register success",
                        "schema": {
                            "$ref": "#/definitions/models.UserAuthDetails"
                        }
                    },
                    "400": {
                        "description": "Failed to register user",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/auth/whoami": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get current user id and name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Authenticated user details",
                "responses": {
                    "200": {
                        "description": "User details",
                        "schema": {
                            "$ref": "#/definitions/models.UserAuthDetails"
                        }
                    },
                    "400": {
                        "description": "Error getting user data",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament": {
            "get": {
                "description": "Get all tournaments",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "All tournaments",
                "responses": {
                    "200": {
                        "description": "Contest bracket",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Tournament"
                            }
                        }
                    },
                    "400": {
                        "description": "Failed to return tournament contest",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/contest/{tournamentId}": {
            "get": {
                "description": "Get tournament contest",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Tournament contest",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tournament id",
                        "name": "tournamentId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "contestType",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Contest bracket",
                        "schema": {
                            "$ref": "#/definitions/models.Bracket"
                        }
                    },
                    "400": {
                        "description": "Failed to return tournament contest",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create new tournament for current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Create new tournament",
                "parameters": [
                    {
                        "description": "Data to create tournament",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateEditTournament"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tournament created",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    },
                    "400": {
                        "description": "Error during tournament creation",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/delete": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete tournaments for current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Delete tournaments",
                "parameters": [
                    {
                        "description": "Data to delete tournaments",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.TournamentIds"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tournaments deleted",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    },
                    "400": {
                        "description": "Error during tournaments deletion",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/delete/{tournamentId}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete tournament for current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Delete tournament",
                "responses": {
                    "200": {
                        "description": "Tournament deleted",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    },
                    "400": {
                        "description": "Error during tournament deletion",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/edit/{tournamentId}": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Edit tournament for current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Edit tournament",
                "parameters": [
                    {
                        "description": "Data to edit tournament",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateEditTournament"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tournament edited",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    },
                    "400": {
                        "description": "Error during tournament edition",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/tiktoks/{tournamentId}": {
            "get": {
                "description": "Get tournament tiktoks",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Tournament tiktoks",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tournament id",
                        "name": "tournamentId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tournament tiktoks",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Tiktok"
                            }
                        }
                    },
                    "400": {
                        "description": "Tournament not found",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/tournament/{tournamentId}": {
            "get": {
                "description": "Get tournament details by its id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tournament"
                ],
                "summary": "Tournament details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tournament id",
                        "name": "tournamentId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tournament",
                        "schema": {
                            "$ref": "#/definitions/models.Tournament"
                        }
                    },
                    "400": {
                        "description": "Tournament not found",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        },
        "/user/tournaments": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create new tournament for current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Create new tournament",
                "responses": {
                    "200": {
                        "description": "Tournaments of user",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    },
                    "400": {
                        "description": "Couldn't get tournaments for specific user",
                        "schema": {
                            "$ref": "#/definitions/controllers.MessageResponseType"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.MessageResponseType": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.AuthInput": {
            "type": "object",
            "required": [
                "name",
                "password"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "models.Bracket": {
            "type": "object",
            "properties": {
                "countMatches": {
                    "type": "integer"
                },
                "rounds": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Round"
                    }
                }
            }
        },
        "models.CreateEditTournament": {
            "type": "object",
            "required": [
                "name",
                "tiktoks"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "size": {
                    "type": "integer",
                    "maximum": 64,
                    "minimum": 4
                },
                "tiktoks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.CreateTiktok"
                    }
                }
            }
        },
        "models.CreateTiktok": {
            "type": "object",
            "required": [
                "url"
            ],
            "properties": {
                "url": {
                    "type": "string"
                }
            }
        },
        "models.Match": {
            "type": "object",
            "properties": {
                "firstOption": {},
                "matchID": {
                    "type": "string"
                },
                "secondOption": {}
            }
        },
        "models.Round": {
            "type": "object",
            "properties": {
                "matches": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Match"
                    }
                },
                "round": {
                    "type": "integer"
                }
            }
        },
        "models.Tiktok": {
            "type": "object",
            "properties": {
                "avgPoints": {
                    "type": "number"
                },
                "timesPlayed": {
                    "type": "integer"
                },
                "tournament": {
                    "$ref": "#/definitions/models.Tournament"
                },
                "tournamentID": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "wins": {
                    "type": "integer"
                }
            }
        },
        "models.Tournament": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "size": {
                    "type": "integer"
                },
                "user": {
                    "$ref": "#/definitions/models.User"
                },
                "userID": {
                    "type": "string"
                }
            }
        },
        "models.TournamentIds": {
            "type": "object",
            "required": [
                "tournamentIds"
            ],
            "properties": {
                "tournamentIds": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "models.UserAuthDetails": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "TikTok arena API",
	Description:      "API for TikTok arena application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
