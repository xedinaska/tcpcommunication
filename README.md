# Description

App is simple TCP server & client. You can find runnable binaries for all platforms in `./bin` folder.
You can build your own app using commands in Makefile (both for client and server)

## Code style

- To fix code format using `gofmt` tool just run `make fmt` command;

## Built With

* Only Golang standard library is used.

## Server usage:

- Run binary using provided binary or `make run` command.

## Client usage:

- Run binary using provided binary or `make run` command.
- Send commands to server in `COMMAND:MESSAGE` format. Use `STOP:` to stop client. Use `SEND:MESSAGE TEXT HERE` to send message from client to server.

## Demonstration
![Alt Text](https://im.ezgif.com/tmp/ezgif-1-796371d274.gif)
