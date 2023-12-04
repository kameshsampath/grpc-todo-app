# Get Started with Data Streaming with gRPC

A demo to show how to get started with building Data Streaming API using [gRPC](https://grpc.dev). We will use [Redpanda](https://redpanda.com) as our streaming platform.

## Prerequisites

To run the demo you need the following tools on your local machine,

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Redpanda CLI](https://docs.redpanda.com/current/get-started/rpk/)
- [gRPC cURL](https://github.com/fullstorydev/grpcurl)
- [Go](https://go.dev)

## Data Streaming Platform Setup

### Start Redpanda Server

```shell
docker compose up -d
```

**IMPORTANT**:

> Wait for the containers to be running. You can use `rpk cluster status` to check the status.

### Create todo-list topic

```shell
rpk topic create todo-list
```

## Start Todo Application Server

Open your Terminal and export the following variables,

```shell
export RPK_BROKERS="127.0.0.1:19092"
# comma separated list of Kafka brokers
export BROKERS="$RPK_BROKERS"
# comma separated list of topics, first topic in the list is default
# producer topic
export TOPICS=todo-list
# allows detailed logging
export ENV=dev
# the port where gRPC Todo app server will be running
export PORT=9090
# The Consumer Group ID
export CONSUMER_GROUP_ID=grpc-todo-app
```

Start the gRPC server,

```shell
go run cmd/server/main.go
```

## Interact With Todo Service

### List Available Services

```shell
grpcurl -plaintext "localhost:$PORT" list
```

Should return an output like,

```text
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
todo.Todo
```

### List Service Methods

```shell
grpcurl -plaintext "localhost:$PORT" list todo.Todo
```

It should return the following methods,

```text
todo.Todo.AddTodo
todo.Todo.TodoList
todo.Todo.UpdateTodo
```

### Start the Streaming Consumer

On a new terminal start the streaming consumer that polls the message on the `todo-list` topic and prints on the console,

```shell
go run cmd/client/main.go
```

**TIP**:

> You can also use
>
> ```shell
> grpcurl -plaintext "localhost:$PORT" todo.Todo/TodoList
> ```

### Add a TODO

```shell
grpcurl -plaintext -d @ "localhost:$PORT" todo.Todo/AddTodo <<EOM
{
  "task": {
    "title": "Finish gRPC Demo README",
    "description": "Complete the README update of the gRPC Data Streaming Demo App.",
    "completed": false
  }
}
EOM
```

Once the task is added the terminal running the client should show the following output,

```text
2023-12-04T21:29:15.823+0530    INFO    client/main.go:43       Task    {"Title": "Finish gRPC Demo README", "Description": "Complete the README update of the gRPC Data Streaming Demo App.", "Completed": false, "Last Updated": "Thursday, 01-Jan-70 00:00:00 UTC", "Partition": 0, "Offset": 1}
```

### Update Todo

```shell
grpcurl -plaintext -d @ "localhost:$PORT" todo.Todo/UpdateTodo <<EOM
{
  "task": {
    "title": "Finish gRPC Demo README",
    "description": "Complete the README update of the gRPC Data Streaming Demo App.",
    "completed": true
  }
}
EOM
```

## Cleanup

```shell
docker compose down
```

Stop the gRPC server and client by hitting `CTRL + C`
