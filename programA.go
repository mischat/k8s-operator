package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    for {
        envValue := os.Getenv("MY_ENV_VAR")
        fmt.Println("Environment Variable Value:", envValue)
        time.Sleep(10 * time.Second)
    }
}

