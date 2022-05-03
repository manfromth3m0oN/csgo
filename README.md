# csgo
A dumb and stupid service made by me Maximillian Stanford-Taylor Jr the IIIerd
----------------------------------------------------------------------------------

{
  "openapi": "3.1.0",
  "info": {
    "title": "Chanshare APIs",
    "description": "API Spec for all services within chanshare\n",
    "version": "0.0.1"
  },
  "paths": {
    "/room": {
      "post": {
        "tags": [
          "Room Manager"
        ],
        "summary": "Start a room",
        "description": "Calling this endpoint should create a room and start its handler in the Room Manager",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "user": {
                    "type": "string"
                  },
                  "playlist": {
                    "type": "array",
                    "items": {
                      "type": "string"
                    }
                  },
                  "room_name": {
                    "type": "string"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "202": {
            "description": "The room has been created"
          },
          "500": {
            "description": "There was an internal error",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "error": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/create-room": {
      "post": {
        "tags": [
          "Room Provisioner"
        ],
        "summary": "Create a new room",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "user": {
                    "type": "string"
                  },
                  "thread_id": {
                    "type": "string"
                  },
                  "board_sn": {
                    "type": "string"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "The room has been created",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "room_name": {
                      "type": "string"
                    },
                    "room_url": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          },
          "500": {
            "description": "There was an internal error",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "error": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/join-room": {
      "post": {
        "tags": [
          "Room Provisioner"
        ],
        "summary": "Join a room",
        "description": "This should send a message to the room manager via NATS",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "user": {
                    "type": "string"
                  },
                  "room-name": {
                    "type": "string"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "The room exists and the user has joined it.\nThe returned URL will be the websocket address for the room\n",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "url": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          },
          "404": {
            "description": "No room was found",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "error": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          },
          "500": {
            "description": "There was an internal error",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "error": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/media/{media_id}": {
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "media_id",
            "required": true,
            "schema": {
              "type": "string"
            }
          }
        ],
        "tags": [
          "Content Service"
        ],
        "summary": "Get a piece of media from the content service",
        "responses": {
          "200": {
            "description": "Returns the found media",
            "content": {
              "video/*": {
                "schema": {
                  "type": "string",
                  "format": "binary"
                }
              }
            }
          },
          "404": {
            "description": "The media isn't in the cache",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "error": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/playlist": {
      "post": {
        "tags": [
          "Content Service"
        ],
        "summary": "Build a playlist and pull all the media",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "board_sn": {
                    "type": "string"
                  },
                  "thread_id": {
                    "type": "string"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successfully started the downloads and built a playlist",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "playlist": {
                      "type": "array",
                      "items": {
                        "type": "string"
                      }
                    }
                  }
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "error": {
                      "type": "string"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/{roomname}": {
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "roomname",
            "required": true,
            "schema": {
              "type": "string"
            }
          }
        ],
        "tags": [
          "Websocket Handler"
        ],
        "summary": "Connect to a room via a ws handler",
        "description": "Should be at the join subdomain e.g. join.chanshare.io/{roomname}\n",
        "responses": {
          "101": {
            "description": "Your connection is being upgraded"
          },
          "404": {
            "description": "The room couldn't be found"
          }
        }
      }
    }
  },
  "components": {}
}
