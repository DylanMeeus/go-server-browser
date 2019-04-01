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
    request(region_europe, "0.0.0.0", "0", "")
    c := make(chan struct{})
    <-c
}

func request(region byte, ip, port, filter string) {
    ipaddr, err := net.ResolveUDPAddr("udp4", addr)
    if err != nil {
        panic(err)
    }
    con, err := net.DialUDP("udp", nil, ipaddr)
    if err != nil {
        panic(err)
    }
    go func(con *net.UDPConn) {
        con.SetReadDeadline(time.Now().Add(10 * time.Second))
        readBytes := make([]byte,1500)
        responseLength, err := con.Read(readBytes)
        if err != nil {
            fmt.Printf("%v\n", err)
        }
        ips := parseResponse(readBytes[:responseLength])
        fmt.Printf("%v\n", ips)
        con.Close()
    }(con)
    _, err = con.Write(compose(0x31, region, ip, port, ""))
    fmt.Println("made request")
    if err != nil {
        panic(err)
    }
}

// return a list of servers based on a response..
func parseResponse(response []byte) ([]string) {
        i := 0
        ips := make([]string,0,len(response) / 6)
        for i < len(response) {
            var fst,snd,thd,fth byte = response[i], response[i+1], response[i+2], response[i+3]
            // read next short?
            var port uint16
            port = uint16(response[i+4]) + uint16(response[i+5])
            ip := fmt.Sprintf("%v.%v.%v.%v:%v\n", fst, snd, thd, fth, port)
            ips = append(ips, ip)
            i += 6
        }
        return ips
}

// compose our message as a series of bytes..
func compose(messageType, region byte, ip, port, filter string) []byte {
    message := make([]byte,0)
    message = append(message, messageType, region)
    // null-terminate our string
    ip += ":" // to connect to port
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
