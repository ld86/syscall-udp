package main

import (
    "os"
    "time"
    "fmt"
    "log"
    "syscall"
    "net"
)

func main() {
    fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)

    if err != nil {
        log.Fatalf("Cannot create socket, %s", err)
    }

    if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
        log.Fatalf("Cannot set SO_REUSEADDR on socket, %s", err)
    }

    udpAddr, err := net.ResolveUDPAddr("udp", ":0")
    if err != nil && udpAddr.IP != nil {
        log.Fatalf("Cannot resolve addr, %s", err)
    }

    if err := syscall.Bind(fd, &syscall.SockaddrInet4{Port: udpAddr.Port}); err != nil {
        log.Fatalf("Cannot bind socket, %s", err)
    }

    file := os.NewFile(uintptr(fd), string(fd))
    conn, err := net.FilePacketConn(file)
    if err != nil {
        log.Fatalf("Cannot create connection from socket, %s", err)
    }

    fmt.Println(conn.LocalAddr())

    if err = file.Close(); err != nil {
        log.Fatalf("Cannot close dup file, %s", err)
    }

    go func() {
        for {
            var buffer [1024]byte
            n, remoteAddr, _ := conn.ReadFrom(buffer[:])

            fmt.Println(string(buffer[:n]))
            fmt.Println(remoteAddr)

            if string(buffer[:n]) == "Ping" {
                conn.WriteTo([]byte("Pong"), remoteAddr)
            }
        }
    }()

    if len(os.Args) > 1 {
        remoteAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
        if err != nil {
            log.Fatalf("Cannot resolve remote addr, %s", err)
        }

        conn.WriteTo([]byte("Ping"), remoteAddr)
    }

    for {
        time.Sleep(time.Second)
    }
}
