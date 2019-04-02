package main

import (
    "encoding/binary"
    "fmt"
    "net"
    "strings"
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
    request(region_europe, "0.0.0.0", "0", "\appid\240", nil)
    c := make(chan struct{})
    <-c
}

func request(region byte, ip, port, filter string, con *net.UDPConn) {
    if con == nil {
        fmt.Println("made new connection..")
        ipaddr, err := net.ResolveUDPAddr("udp4", addr)
        if err != nil {
            panic(err)
        }
        con, err = net.DialUDP("udp", nil, ipaddr)
        if err != nil {
            panic(err)
        }
    }
    go func(con *net.UDPConn) {
        readBytes := make([]byte,3000)
        responseLength, err := con.Read(readBytes)
        if err != nil {
            fmt.Printf("%v\n", err)
        }
        ips := parseResponse(readBytes[:responseLength])
        for _,i := range ips {
            if i == "0.0.0.0:0" {
                print("YAY!!")
                return
            }
        }
        lastIp := ips[len(ips)-1]
        fmt.Println(lastIp)
        if lastIp != "0.0.0.0:0" {
            parts := strings.Split(lastIp, ":")
            request(region, parts[0], parts[1], filter, con) 
        } else {
            // last ip was received so close the connection
            fmt.Println("closing connection..")
            con.Close()
        }
    }(con)
    _, err := con.Write(compose(0x31, region, ip, port, filter)) 
    fmt.Println("made request")
    if err != nil {
        panic(err)
    }
}

// return a list of servers based on a response..
func parseResponse(response []byte) ([]string) {
       i := 6
       ips := make([]string,0)
        for i < len(response) {
            var fst,snd,thd,fth byte = response[i], response[i+1], response[i+2], response[i+3]
            // read next short?
            s := []byte{response[i+4], response[i+5]}
            port := binary.BigEndian.Uint16(s)
            ip := fmt.Sprintf("%v.%v.%v.%v:%v", fst, snd, thd, fth, port)
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
