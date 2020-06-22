package handler

import (
    "os"
    "log"
)

func JsonWrite(data []byte) {
    fp, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0755)
    if err != nil {
        log.Fatal(err)
    }
    defer fp.Close()
    _, err = fp.Write(data)
    if err != nil {
        log.Fatal(err)
    }
}

