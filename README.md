# kval
A simple in-memory key-value store with time duration. Kval is essentially an in-memory dictionary with concurrent read/write access. Kval can be run as a binary or embedded in a go program. It has a built in RPC and HTTP server.

This project serves as a way to learn key-value databases by building a simple one. 

Features:
* Concurrency safe
* Sharded to minimize lock contention


## Installation
Assuming go enviorment variables are set up correctly, run:

```
go get github.com/andresoro/kval
```

## Usage

### Standalone

In one terminal start the server. Kval will start an RPC server and an http server.
```
kval server
```
In another terminal start the kval cli. Commands in cli are very simple: add, get, delete, exit

```
kval cli
>>> add key val
Key-Value successfully added
>>> get key
val
```

### In a go program

``` go
import "github.com/andresoro/kval/kval"
import "time"

type MyStruct struct {
    Example string
    Number int
}

func main() {
    kval := kval.New(4, time.Minute)
    x := MyStruct{
        Example: "string",
        Number: 4,
    }

    kval.Add("key", x)

    //CONTINUE ON 
}

```



Todo:

* Documentation
* Config
* Add a max item size and/or max cache size
* Add file persistence for easy restarting
* ~~Logging~~
* ~~Add sharding~~
* ~~Add item expiration (Use FIFO queue)~~
* ~~RPC Server/Cli~~
* ~~Logging~~


