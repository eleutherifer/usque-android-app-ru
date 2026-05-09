package main

import (
  "bufio"
  "context"
  "crypto/tls"
  "fmt"
  "io"
  "net"
  "net/http"
  "os"
  "strings"
  "time"

  utls "github.com/refraction-networking/utls"
)

func rawUTLS(ip string, port string, id utls.ClientHelloID) string {
  ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second); defer cancel()
  d := net.Dialer{Timeout: 3500*time.Millisecond}
  raw, err := d.DialContext(ctx, "tcp", net.JoinHostPort(ip, port)); if err != nil { return "dial " + err.Error() }
  defer raw.Close(); _ = raw.SetDeadline(time.Now().Add(8*time.Second))
  c := utls.UClient(raw, &utls.Config{ServerName:"speed.cloudflare.com", InsecureSkipVerify:true, MinVersion:utls.VersionTLS12}, id)
  if err := c.Handshake(); err != nil { return "tls " + err.Error() }
  req := "GET /__down?bytes=1000000 HTTP/1.1\r\n"+
    "Host: speed.cloudflare.com\r\n"+
    "User-Agent: Mozilla/5.0\r\n"+
    "Accept: */*\r\n"+
    "Accept-Encoding: identity\r\n"+
    "Connection: close\r\n\r\n"
  if _, err := c.Write([]byte(req)); err != nil { return "write " + err.Error() }
  br := bufio.NewReader(c)
  st, err := br.ReadString('\n'); if err != nil { return "status " + err.Error() }
  for { l, er := br.ReadString('\n'); if er != nil { return "hdr " + er.Error() }; if l == "\r\n" { break } }
  n, _ := io.CopyN(io.Discard, br, 1024)
  return strings.TrimSpace(st) + fmt.Sprintf(" bytes=%d", n)
}

func stdHTTP(ip, port string) string {
  d := net.Dialer{Timeout:3500*time.Millisecond}
  tr := &http.Transport{ForceAttemptHTTP2:false, DisableKeepAlives:true, TLSClientConfig:&tls.Config{ServerName:"speed.cloudflare.com", InsecureSkipVerify:true, NextProtos:[]string{"http/1.1"}}, DialContext: func(ctx context.Context, network, addr string)(net.Conn,error){ return d.DialContext(ctx, network, net.JoinHostPort(ip,port)) }}
  c := &http.Client{Transport:tr, Timeout:8*time.Second}
  r, err := http.NewRequest("GET", "https://speed.cloudflare.com/__down?bytes=1000000", nil); if err != nil { return err.Error() }
  r.Host="speed.cloudflare.com"; r.Header.Set("User-Agent","Mozilla/5.0"); r.Header.Set("Accept-Encoding","identity")
  resp, err := c.Do(r); if err != nil { return "err " + err.Error() }
  defer resp.Body.Close(); n,_:=io.CopyN(io.Discard, resp.Body, 1024)
  return resp.Status + fmt.Sprintf(" bytes=%d", n)
}

func main(){
  nodes := []string{"23.141.52.169:443", "103.46.142.56:443", "176.98.181.71:443", "154.219.103.79:443"}
  if len(os.Args)>1 { nodes = os.Args[1:] }
  ids := []struct{name string; id utls.ClientHelloID}{
    {"ChromeAuto", utls.HelloChrome_Auto},
    {"Chrome120", utls.HelloChrome_120},
    {"Chrome102", utls.HelloChrome_102},
    {"Firefox120", utls.HelloFirefox_120},
    {"IOSAuto", utls.HelloIOS_Auto},
    {"Randomized", utls.HelloRandomized},
    {"RandomizedALPN", utls.HelloRandomizedALPN},
    {"RandomizedNoALPN", utls.HelloRandomizedNoALPN},
  }
  for _, node := range nodes {
    parts := strings.Split(node, ":"); ip, port := parts[0], parts[1]
    fmt.Println("\n##", node)
    fmt.Println("std", stdHTTP(ip,port))
    for _, x := range ids { fmt.Println(x.name, rawUTLS(ip,port,x.id)) }
  }
}
