package speedgo

import (
    "context"
    "crypto/tls"
    "fmt"
    "io"
    "net"
    "net/http"
    "net/http/httptrace"
    "net/url"
    "strconv"
    "strings"
    "time"
)

const defaultSpeedURL = "https://speed.cloudflare.com/__down?bytes=10000000"

func TestSpeed(host string, port int, timeoutMs int) string {
    return TestSpeedURL(host, port, defaultSpeedURL, timeoutMs)
}

func TestSpeedURL(host string, port int, speedURL string, timeoutMs int) string {
    if timeoutMs <= 0 { timeoutMs = 9000 }
    if strings.TrimSpace(speedURL) == "" { speedURL = defaultSpeedURL }
    if !strings.HasPrefix(speedURL, "http://") && !strings.HasPrefix(speedURL, "https://") { speedURL = "https://" + speedURL }
    u, err := url.Parse(speedURL)
    if err != nil { return "ERR bad url: " + err.Error() }
    targetHost := u.Hostname()
    if targetHost == "" { return "ERR bad url: empty host" }

    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMs)*time.Millisecond)
    defer cancel()
    forcedAddr := net.JoinHostPort(host, strconv.Itoa(port))
    dialer := &net.Dialer{Timeout: 3500 * time.Millisecond, KeepAlive: 30 * time.Second}

    tlsInfo := "tls=?"
    trace := &httptrace.ClientTrace{
        TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
            if err != nil { tlsInfo = "tls_err=" + err.Error(); return }
            tlsInfo = fmt.Sprintf("tls=v%x alpn=%s cipher=0x%x server=%s", cs.Version, cs.NegotiatedProtocol, cs.CipherSuite, cs.ServerName)
        },
    }
    ctx = httptrace.WithClientTrace(ctx, trace)

    tr := &http.Transport{
        Proxy:                 nil,
        ForceAttemptHTTP2:     false,
        DisableCompression:    true,
        DisableKeepAlives:     true,
        TLSHandshakeTimeout:   3500 * time.Millisecond,
        ResponseHeaderTimeout: 3500 * time.Millisecond,
        ExpectContinueTimeout: 1 * time.Second,
        TLSClientConfig: &tls.Config{
            ServerName:         targetHost,
            InsecureSkipVerify: true,
            MinVersion:         tls.VersionTLS12,
            NextProtos:         []string{"http/1.1"},
        },
        DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
            return dialer.DialContext(ctx, network, forcedAddr)
        },
    }
    defer tr.CloseIdleConnections()
    client := &http.Client{Transport: tr, Timeout: time.Duration(timeoutMs) * time.Millisecond}

    // curl 等价形态：Request URL/Host/SNI 都保持 speed.cloudflare.com，
    // 只在 DialContext 里把底层 TCP 连接改到候选 IP:port。
    // 非 443 端口必须写进 URL Host，否则 net/http 会按 443 建连接目标。
    reqURL := *u
    defaultPort := "443"
    if reqURL.Scheme == "http" { defaultPort = "80" }
    if portStr := strconv.Itoa(port); portStr != defaultPort {
        reqURL.Host = net.JoinHostPort(targetHost, portStr)
    } else {
        reqURL.Host = targetHost
    }
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
    if err != nil { return "ERR request: " + err.Error() }
    req.Host = targetHost
    req.Header.Set("User-Agent", "Mozilla/5.0")
    req.Header.Set("Accept", "*/*")
    req.Header.Set("Accept-Encoding", "identity")
    req.Header.Set("Connection", "close")

    start := time.Now()
    resp, err := client.Do(req)
    if err != nil { return "ERR http: " + err.Error() + " [" + tlsInfo + "]" }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        small, _ := io.ReadAll(io.LimitReader(resp.Body, 160))
        body := strings.Join(strings.Fields(string(small)), " ")
        if len(body) > 80 { body = body[:80] }
        return fmt.Sprintf("ERR HTTP %d [%s] server=%s cf-ray=%s body=%q", resp.StatusCode, tlsInfo, resp.Header.Get("Server"), resp.Header.Get("Cf-Ray"), body)
    }
    buf := make([]byte, 32*1024)
    var total int64
    deadline := time.Now().Add(time.Duration(timeoutMs) * time.Millisecond)
    for total < 32_000_000 && time.Now().Before(deadline) {
        n, er := resp.Body.Read(buf)
        if n > 0 { total += int64(n) }
        if er == io.EOF { break }
        if er != nil {
            if total == 0 { return "ERR read: " + er.Error() + " [" + tlsInfo + "]" }
            break
        }
    }
    sec := time.Since(start).Seconds()
    if total <= 0 || sec <= 0 { return "ERR 0 bytes [" + tlsInfo + "]" }
    mbps := (float64(total) / 1024.0 / 1024.0) / sec
    return fmt.Sprintf("OK %.2f [std %s]", mbps, tlsInfo)
}
