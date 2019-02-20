# kval
A simple in-memory key-value store with time duration. Kval is essentially an in-memory dictionary with concurrent read/write access. Kval can be run as a binary or embedded in a go program. It has a built in RPC and HTTP server.

This project serves as a way to learn in memory key-value databases by building a simple one. 

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

All http requests can be sent to same url, localhost:PORT/kval/{key}. Where key is desired key. PUT requests need the value in the body of the requests. 

### In a go program

``` go
import (
    "github.com/andresoro/kval/kval"
    "github.com/andresoro/kval/server"
    "time"
    "bytes"
)

type MyStruct struct {
    Example string
    Number int
}

func main() {
    // new store with 4 shards and a one minute eviction policy
    kval, _ := kval.New(4, time.Minute)
    x := MyStruct{
        Example: "string",
        Number: 4,
    }

    buf := new(bytes.Buffer)
    gob.NewEncoder(buf).Encode(x)

    // add to key-value store
    err := kval.Add("key", buf.Bytes())
    if err != nil {
        // do something
    }

    // get value from store
    val, err := kval.Get("key")
    if err != nil {
        //do something
    }

    // delete returns value and an error
    val, err := kval.Delete("key")
    if err != nil {
        // do something
    }

    //start an http server with previously defined store
    server := server.NewHTTP(":8080", kval)
    server.Start()

}

```





Todo:

* Add a max item size/ and or max cache size
* ~~Logging~~
* ~~Add sharding~~
* ~~Add item expiration (Use FIFO queue)~~
* ~~RPC Server/Cli~~
* ~~Logging~~


