package main

import (
    "github.com/mvillalba/openexchangerates/oxr"
    "os"
    "fmt"
)

func main() {
    // Check args
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "app_id")
        os.Exit(64) // EX_USAGE
    }

    // Init ApiClient
    client := oxr.New(os.Args[1])

    // Fun stuff
    currencies(client)
    latest(client)
    historical(client)
}

func currencies(client *oxr.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List all available currencies.")
    fmt.Println("=======================================")

    curr, err := client.Currencies()
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    for k := range curr {
        fmt.Println(k, curr[k])
    }
}

func latest(client *oxr.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List latest exchange rates.")
    fmt.Println("=======================================")

    r, err := client.Latest()
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }
    fmt.Println("Disclaimer:", r.Disclaimer)
    fmt.Println("License:", r.License)
    fmt.Println("Timestamp:", r.Timestamp)
    fmt.Println("Base:", r.Base)
    for k := range r.Rates {
        fmt.Println("USD/" + k, r.Rates[k])
    }
}

func historical(client *oxr.ApiClient) {
    fmt.Println()
    fmt.Println("=======================================")
    fmt.Println("List latest exchange rates as they were")
    fmt.Println("on 2014-01-01.")
    fmt.Println("=======================================")

    r, err := client.Historical("2014-01-01")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }
    for k := range r.Rates {
        fmt.Println("USD/" + k, r.Rates[k])
    }
}
