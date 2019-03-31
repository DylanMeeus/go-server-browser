package main

import (
    "fmt"
    "net"
    "time"
)

// region
const (
    region_us_east = 0x00
    region_us_west = 0x01
    region_south_america =  0x02
    region_europe = 0x03
)

const (
    addr = "hl2master.steampowered.com:27011"
    port = "27011"
)

func main() {
    ipaddr, err := net.ResolveUDPAddr("udp4", addr)
    if err != nil {
        panic(err)
    }
    con, err := net.DialUDP("udp", nil, ipaddr)
    if err != nil {
        panic(err)
    }
    go func(con *net.UDPConn) {
        fmt.Println("started reading..")
        con.SetReadDeadline(time.Now().Add(10 * time.Second))
        readBytes := make([]byte,1500)
        res, err := con.Read(readBytes)
        fmt.Println("read")
        if err != nil {
            fmt.Printf("%v\n", err)
        }
        fmt.Printf("%v\n", res)
        fmt.Printf("%v\n", readBytes)
        con.Close()
    }(con)
    _, err = con.Write(compose(0x31, region_us_east, "0.0.0.0:", "0", ""))
    fmt.Println("written data..")
    if err != nil {
        panic(err)
    }
    c := make(chan struct{})
    <-c
}

// compose our message as a series of bytes..
func compose(messageType, region byte, ip, port, filter string) []byte {
    message := make([]byte,0)
    message = append(message, messageType, region)
    // null-terminate our string
    port += "\000"
    filter += "\000"
    for _,b := range []byte(ip) {
        message = append(message, b)
    }
    for _,b := range []byte(port) {
        message = append(message, b)
    }
    for _,b := range []byte(filter) {
        message = append(message, b)

    }
    return message
}
