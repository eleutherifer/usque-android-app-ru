package main

import (
  "fmt"
  "speedgo"
)
func main(){
  nodes := []struct{ip string; port int}{
    {"68.64.182.39",443},{"23.175.201.2",8443},{"219.76.13.169",443},{"103.46.142.56",443},{"202.61.72.141",443},{"47.79.20.130",443},{"43.254.164.50",8443},
  }
  ok:=0
  for _, n := range nodes {
    r:=speedgo.TestSpeed(n.ip,n.port,9000)
    fmt.Printf("%s:%d => %s\n", n.ip, n.port, r)
    if len(r)>=3 && r[:3]=="OK " { ok++ }
  }
  if ok == 0 { panic("no successful speed test") }
  fmt.Printf("OK_COUNT=%d/%d\n", ok, len(nodes))
}
