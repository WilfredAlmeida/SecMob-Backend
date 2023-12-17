# SecMob

SecMob facilitates executing predefined commands on a server via a REST API. This repo is the backend that runs on the server

The backend runs as a REST API. Executable commands are defined in a file, and the system parses them and facilitates their execution.

## Basic Idea of SecMob

Executable commands are predefined on the server in a format (JSON in this case). These commands are parsed by the backend system and sent over to the client via a REST API. The client receives randomly unique IDs of the commands and requests execution of the commands by sending back the ID to the server. The server executed the command defined on its end and returns the result to the client. The communication between the client and server is encrypted using AES-256-GCM and optional SSL.

* [Architecture](#architecture)

* [Authentication](#authentication)

    * [Additional authentication feature](#additional-authentication-feature)

* [AES Key Exchange](#aes-key-exchange)

* [Defining Commands](#defining-commands)

* [API Reference](#api-reference)

* [Directory Structure](#directory-structure)

* [Running Locally](#running-locally)

* [Building](#building)

* [Contributing](#contributing)

* [Taking to Production](#taking-to-production)


## Architecture

Following is the architecture of the system

![](https://cdn.hashnode.com/res/hashnode/image/upload/v1671626896205/CqAz8qCmK.png)

**Steps Explained**

1. Requests are received

2. Requests get authenticated

3. If authentication fails, an unauthorized response is sent to the client

4. Authorized requests are processed by respective handlers


Steps 5 - 9 are for `Get Commands` API endpoint

1. The handler fetches available commands from the commands parser and returns them

2. The commands parser module reads the commands from the file system and parses them in the desired format

3. The commands are defined in a file in a format that gets parsed

4. The commands parser module gives the parsed commands to the handler

5. The handler returns the commands to be sent as the response


Steps 10 - 15 are for `Run Commands` API endpoint

1. The handler executes the specified command and returns the result

2. The handler gives the command ID of the command to be executed to the command execution module

3. The command execution module asks for the command to execute from the commands parser module

4. The commands parser module returns the command to be executed

5. The command execution module executes the command

6. The execution output is sent to the handler

7. The handler returns the execution result to be sent as the response


Steps 17 - 19 are generalized for the whole system

1. The response module handles sending of responses

2. The response body is prepared as desired and encrypted

3. The response is sent to the client


## Authentication

The API body exchanged between the client and server is encrypted using AES-256-GCM. If the server can decrypt the request body sent by the client, it means it uses the correct key, and thus the client is authenticated.

A session is not maintained between the client and server, every incoming request's body gets decrypted and authentication is performed.

### Additional authentication feature

Additionally, a unique device ID of the client device can be stored on the server. If the incoming requests have the body then they'll be authenticated. This feature is yet to be added because the client-side Flutter application does not have a reliable solution to generate device-specific IDs on Android.

## AES Key Exchange

The AES key needs to be stored on the server in a file. Key generation and exchange are not facilitated by SecMob. The user must generate the key and store it on the server as well as on the mobile application. On the mobile application, the key is accepted by scanning a QR code.

## Defining Commands

The commands must be defined in the file as specified further. The user should make sure that the defined command can be executed from the application and OS user's scope w.r.t. authorization and permissions.

The commands file is auto-reloaded on save and a restart is not needed to deploy modified/new commands.

## API Reference

API responses are base64 encoded, AES-256-GCM encoded strings.

1. `/getCommands`: HTTP Post request to get the list of available commands from the server.


**Request Body**: None

**Response**:

Following is the JSON description

```json
{
   "status":1,
   "message":"Data Fetched Successfully",
   "payload":[
      {
         "id":"UfKvDP",
         "title":"Get Docker Images"
      },
      {
         "id":"ttLNpU",
         "title":"Echo Hi"
      },
      {
         "id":"krgDkm",
         "title":"Git Status"
      }
   ]
}
```

`status`: Indicates the status of command execution. Values

* `0`: Command execution failed with an error

* `1`: Command execution successful


`message`: Informative message regarding the operation

`payload`: Response payload. Is always a list of objects. The objects contain response data

Command Object:

```json
{
    "id":"krgDkm",
    "title":"Git Status"
}
```

`id`: Randomly generated ID for commands by the server. Changes when the server restarts or commands are refreshed.

`title`: Descriptive title of the command

Response when no commands are defined

```json
{
   "message":"No Commands Found",
   "payload":null,
   "status":1
}
```

1. `/runCommand`: HTTP Post request to send execution request to the server


Request Body:

```json
{
    "commandId":"uXrXrS"
}
```

The server will parse the command with the provided ID and send its execution response

Response:

1. Command execution successful

   HTTP Status: `200`

    ```json
    {
       "message":"Command Executed Successfully",
       "payload":[
          {
             "output":"On branch main nothing to commit, working tree clean"
          }
       ],
       "status":1
    }
    ```

   Payload Object: The payload is always a list of objects. The object contains the execution result as follows

    ```json
    {
        "output":"On branch main nothing to commit, working tree clean"
    }
    ```

   `output`: A string of command execution responses.

2. Command execution failed

   HTTP Status: `200`

    ```json
    {
       "message":"Command Execution Failed",
       "payload":[
          {
             "output":"exec: \"echo\": executable file not found in %PATH%"
          }
       ],
       "status":0
    }
    ```

   Note the `status` here.

3. Invalid Body

   HTTP Status: `206`

    ```json
    {
        "message": "Invalid Body",
        "payload": null,
        "status": 0
    }
    ```


## Directory Structure

```bash
.
├── apiRoutes
│   └── api.go
├── commands
│   └── commands.json
├── crypto
│   ├── aes_key_loader.go
│   ├── decryption.go
│   └── encryption.go
├── go.mod
├── go.sum
├── keys
│   └── aesKey
├── main.go
├── qrcode.png
└── utils
    ├── constants.go
    ├── file_watcher_routine.go
    ├── globals.go
    ├── logger.go
    ├── random_string_generator.go
    └── scripts
        └── generateQr.go
```

`.env`: Consists of environment variables. The following are the ones needed:

* `PORT`: Port on which the server listens

* `GIN_MODE`: Operation mode for Gin: Possible values are as follows:

    * `debug`

    * `release`


`go.mod, go.sum`: Module files for Go. Nothing to do here except modules and dependencies

`main.go`: The starting point of the application

`SecMob-logs.log`: Log file. The path can be changed in `main.go`

`commands/commands.json`: The JSON file where commands are defined. Following is the schema for a command

```json
{
    "title": "Get Docker Images",
    "command": "docker image ls"
}
```

Commands have IDs auto-generated when the file is parsed. Restart the application to refresh the commands.

`crypto/encryption.go`: Has functions to encrypt data

`crypto/decryption.go`: Has functions to decrypt data

`crypto/aes_key_loader.go`: Loads AES key from the file into a global variable

`keys/aesKey`: The AES key file. It's a plaintext file that has the AES key

EXTREMELY IMPORTANT: ADD THE KEY FILE TO `.gitignore` ON YOUR SYSTEM AND KEEP IT AWAY FROM VCS. IT IS SHARED HERE AS AN EXAMPLE. IT IS RECOMMENDED TO CHANGE THE KEY FILE PATH IN `main.go` AND STORE YOUR KEYS SOMEPLACE SECURE.

`apiRoutes/command.go`: Has API handlers

`utils`: This directory consists of utility functions

`utils/constants.go`: Has constant variables

`utils/file_watcher_routine.go`: Has a function that watches a given file for changes. Used to watch `commands/commands.json`

`utils/random_string_generator.go`: Generates a random string of a given length

`utils/globals.go`: Has global variables

`utils/logger.go`: Logger and its functions  

`utils/scripts/generateQr.go`: Generates QR from given user id and AES key. Share this QR to the user

## Running Locally

1. Clone the repo

2. Set environment variables. Set it to `debug` if you're running for the first time or are in development

3. Change file paths in `main.go` if desired

4. Run `go run main.go`

5. If it's `debug` mode then you'll see some logs printed on terminal

6. That's it, you now have it running


## Building

This [tutorial](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04) by DigitalOcean is a great resource for building the application. Based on your circumstances, build accordingly.

On Windows, in git bash, this command works for me:

```bash
env GOOS=windows GOARCH=amd64 go build
```
