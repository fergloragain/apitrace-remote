# apitrace-remote
Traces graphical applications remotely with apitrace, intended to be used in conjunction with the [apitrace-remote-viewer](/fergloragain/apitrace-remote-viewer)

- [Building](#building)
- [Running](#running)
- [User flow](#user-flow)
  * [Create an app](#create-an-app)
  * [Trace the app](#trace-the-app)
  * [Get the dump](#get-the-dump)
  * [View the GL state](#view-the-gl-state)
- [Endpoints](#endpoints)
  * [Apps](#apps)
    + [GET `/apps`](#get-apps)
      - [Request](#request)
      - [Response](#response)
    + [POST `/apps/:name`](#post-appsname)
      - [Request](#request-1)
      - [Response](#response-1)
    + [GET `/apps/:name`](#get-appsname)
      - [Request](#request-2)
      - [Response](#response-2)
    + [PUT `/apps/:name`](#put-appsname)
      - [Request](#request-3)
      - [Response](#response-3)
  * [Traces](#traces)
    + [GET `/traces`](#get-traces)
      - [Request](#request-4)
      - [Response](#response-4)
    + [POST `/traces/:name`](#post-tracesname)
      - [Request](#request-5)
      - [Response](#response-5)
    + [GET `/traces/:name`](#get-tracesname)
      - [Request](#request-6)
      - [Response](#response-6)
- [Todo](#todo)

## Building

go build -o main cmd/server/*go

## Running 

./main 

## User flow

- Create an app
- Trace the app
- Get the dump
- View the GL state 

### Create an app

Add the details for an application to be traced. Details include:

- Name
- Description
- Git URL
- User
- Private key path
- Branch 
- Build script
- Executable
- Timeout
- apitrace location
- glretrace location

### Trace the app

Clones the git repo to a folder on disk, then runs the build script to produce the executable. Once built, the following command is ran against the executable:

```bash
apitrace trace myapp
```

With the trace file written on disk, apitrace is used to dump the per-frame GL calls with the following command:

```bash
apitrace dump myapp.trace
```

### Get the dump

Retrieve the list of GL calls made for a specific frame of the app, beginning with frame 0. Clicking on a particular GL call will trigger a glretrace run to capture the GL state for that particular call:

```bash
glretrace -D=12345 myapp.trace
```

### View the GL state

Following the call to `glretrace`, the colour, depth, and stencil buffers are viewable, as well as the GL state, including uniforms, shaders, buffers, etc

## Endpoints 

### Apps

#### GET `/apps`

Retrieves a list of the applications stored in the database

##### Request 

```bash
curl -X GET http://localhost:8080/apps
```

##### Response 

```json
[{"id":"hellmouthxyz","name":"hellmouthxyz","description":"jjkjklj","branch":"apitraceremote","traces":3}]
```

#### POST `/apps/:name`

Creates a new application in the database 

##### Request 

```bash
curl -X POST -d '{"description":"jjkjklj","url":"https://github.com/fergloragain/hellmouthxyz.git","executable":"main","apiTrace":"/Users/hellmouthxyz/development/apitrace/build/apitrace","retrace":"/Users/hellmouthxyz/development/apitrace/build/glretrace","timeout":2,"user":"","privateKey":"","buildScript":"build.sh","branch":"apitraceremote","dumpImages":true}' http://localhost:8080/apps/hellmouthxyztest
```

##### Response 

```json
{"id":"hellmouthxyztest-6","name":"hellmouthxyztest","description":"jjkjklj","url":"https://github.com/fergloragain/hellmouthxyz.git","executable":"main","apiTrace":"/Users/hellmouthxyz/development/apitrace/build/apitrace","retrace":"/Users/hellmouthxyz/development/apitrace/build/glretrace","timeout":2,"user":"","privateKey":"","buildScript":"build.sh","active":false,"branch":"apitraceremote","traces":[],"dumpImages":true}
```

#### GET `/apps/:name`

Gets the details for the `:name` app in the database

##### Request 

```bash
curl -X GET http://localhost:8080/apps/hellmouthxyz-6
```

##### Response 

```json
{"id":"hellmouthxyz-6","name":"hellmouthxyz","description":"jjkjklj","url":"https://github.com/fergloragain/hellmouthxyz.git","executable":"main","apiTrace":"/Users/hellmouthxyz/development/apitrace/build/apitrace","retrace":"/Users/hellmouthxyz/development/apitrace/build/glretrace","timeout":2,"user":"","privateKey":"","buildScript":"build.sh","active":false,"branch":"apitraceremote","traces":["hellmouthxyz-trace-1","hellmouthxyz-trace-2","hellmouthxyz-trace-3"],"dumpImages":true}
``` 

#### PUT `/apps/:name`

Updates an existing application in the database 

##### Request 

```bash
curl -X PUT -d '{"description":"abc","url":"https://github.com/fergloragain/hellmouthxyz.git","executable":"main","apiTrace":"/Users/hellmouthxyz/development/apitrace/build/apitrace","retrace":"/Users/hellmouthxyz/development/apitrace/build/glretrace","timeout":2,"user":"","privateKey":"","buildScript":"build.sh","branch":"apitraceremote","dumpImages":true}' http://localhost:8080/apps/hellmouthxyztest-6
```

##### Response 

```json
{"id":"hellmouthxyztest-6","name":"hellmouthxyztest","description":"abc","url":"https://github.com/fergloragain/hellmouthxyz.git","executable":"main","apiTrace":"/Users/hellmouthxyz/development/apitrace/build/apitrace","retrace":"/Users/hellmouthxyz/development/apitrace/build/glretrace","timeout":2,"user":"","privateKey":"","buildScript":"build.sh","active":false,"branch":"apitraceremote","traces":[],"dumpImages":true}
```

### Traces

#### GET `/traces`

Retrieves a list of the traces stored in the database

##### Request 

```bash
curl -X GET http://localhost:8080/traces
```

##### Response 

```json
["hellmouthxyz-1-trace","hellmouthxyz-1-trace-1"]
```

#### POST `/traces/:name`

Creates a new trace in the database for the `:name` app

##### Request 

```bash
curl -X POST http://localhost:8080/traces/hellmouthxyztest
```

##### Response 

```json
{"id":"hellmouthxyztest-trace","appID":"hellmouthxyztest","name":"hellmouthxyztest-trace","status":"Pending","buildStdout":"","buildStderr":"","traceStdout":"","traceStderr":"","cloneStdout":"","cloneStderr":"","dumpStderr":"","targetDirectory":"","numberOfFrames":0,"retraces":[],"traceFile":""}
```

#### GET `/traces/:name`

Gets the details for the `:name` trace in the database

##### Request 

```bash
curl -X GET http://localhost:8080/traces/hellmouthxyz-23-trace
```

##### Response 

```json
{"id":"hellmouthxyz-23-trace","appID":"hellmouthxyz-23","name":"hellmouthxyz-23-trace","status":"Pending","buildStdout":"","buildStderr":"","traceStdout":"","traceStderr":"","cloneStdout":"","cloneStderr":"","dumpStderr":"","targetDirectory":"","numberOfFrames":0,"retraces":[],"traceFile":""}
``` 

## Todo 

- [ ] Make the project `go get` friendly
- [ ] Add config file for default values for new apps
- [ ] Add Docker files and a docker-compose.yaml
- [ ] Fix up deletion so that deletion of a project triggers deletion of:
    - [ ] All DB data
        - [ ] Retrace data
        - [ ] Dump data
        - [ ] App data
    - [ ] All on disk data
        - [ ] Project folder
        - [ ] Trace directories
- [ ] Add pagination for viewing the dump; i.e. if a single frame makes 10,000 calls, then return the first 100 calls, and retrieve the next set when the page is scrolled to the bottom
- [ ] Trigger the `glretrace` operations asynchronously, and have the client application poll the status repeatedly
- [ ] Add the a profiling call for `glretrace`
- [ ] Add logic so that instead of re-cloning the source code every time, a git pull is performed and a rebuild triggered, unless explicitly stated otherwise
- [ ] Move the executable into a separate trace folder 
- [ ] Add disk statistics to the data set returned, so we know how much disk space a particular application/trace is using
- [ ] Rework and streamline the application triggering process, as well as stdout and stderr capture


