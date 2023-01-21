package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Data struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // Receive data from first API
    res, err := http.Get("http://firstapi.com/data")
    if err != nil {
        panic(err)
    }
    defer res.Body.Close()

    // Read response body and unmarshal into struct
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic(err)
    }

    var data Data
    json.Unmarshal(body, &data)

    // Send data to second API
    url := "http://secondapi.com/data"
    jsonData, _ := json.Marshal(data)
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // Print response from second API
    fmt.Println("Response from second API:", resp.Status)
}
