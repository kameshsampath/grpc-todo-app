# Get Started with Data Streaming with gRPC

A demo to show how to get started with building Data Streaming API using [gRPC](https://grpc.dev). We will use [Redpanda](https://redpanda.com) as our streaming platform.

## Prerequisites

To run the demo you need the following tools on your local machine,

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Redpanda CLI](https://docs.redpanda.com/current/get-started/rpk/)
- [gRPC cURL](https://github.com/fullstorydev/grpcurl)

## Data Streaming Platform Setup

The entire demo is containerized and all the applications could be started using the Docker compose,

```shell
docker compose up -d
```

The docker compose starts the following services,

- `redpanda-0` - A single node Redpanda server
- `console` - The Redpanda console
- `todo-app-server` - The Todo Application gRPC server
- `todo-list` - The Todo Application client that receives the streaming messages from the gRPC `todo-app-server`

The `Todo` application runs with the following environment variables defaults, please change them as needed if you deploy without these defaults.

### Todo gRPC Server

```shell
# gRPC service port
PORT=9090
# Redpanda Brokers
BROKERS=redpanda-0:9092
# Topic to store the Todo
TOPICS=todo-list
# Running environment, typically used for grpcurl
ENV=dev
# The consumer group used while consuming messages
CONSUMER_GROUP_ID=grpc-todo-app
```

### Todo gRPC Client (Todo List)

```shell
SERVICE_ADDRESS=todo-app-server:9090
```

**NOTE**:

> The individual application binaries for Todo App gRPC Server and Todo App Client are available on the [application repo](https://github.com/kameshsampath/grpc-todo-app/releases). You can download them and run the application individually.

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
```

### View List of Todo

On a new terminal run the following command to view the list of Todos added by earlier steps,

```shell
docker compose logs -f todo-list
```

The output should be something like,

```shell
todo-list  | 2023-12-05T05:20:19.235Z   INFO    client/main.go:44       Task    {"Title": "Finish gRPC Demo README", "Description": "Complete the README update of the gRPC Data Streaming Demo App.", "Completed": false, "Last Updated": "Thursday, 01-Jan-70 00:00:00 UTC", "Partition": 0, "Offset": 0}
todo-list  | 2023-12-05T05:20:19.236Z   INFO    client/main.go:44       Task    {"Title": "Finish gRPC Demo README", "Description": "Complete the README update of the gRPC Data Streaming Demo App.", "Completed": false, "Last Updated": "Thursday, 01-Jan-70 00:00:00 UTC", "Partition": 0, "Offset": 1}
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

Once the task is added the terminal running the client should show an output similar to,

```json
{
  "Title": "Finish gRPC Demo README",
  "Description": "Complete the README update of the gRPC Data Streaming Demo App.",
  "Completed": false,
  "Last Updated": "Thursday, 01-Jan-70 00:00:00 UTC",
  "Partition": 0,
  "Offset": 1
}
```

## References

The demo uses the [protobuf](https://protobuf.dev) definitions from <https://github.com/kameshsampath/demo-protos>

## Cleanup

```shell
docker compose down
```
