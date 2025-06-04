package loadbal

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIp(r *http.Request) string {
    // Check X-Forwarded-For header first (in case of proxy chain)
    if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
        ips := strings.Split(xff, ",")
        return strings.TrimSpace(ips[0])
    }
    
    // Check X-Real-IP header
    if xri := r.Header.Get("X-Real-IP"); xri != "" {
        return strings.TrimSpace(xri)
    }
    
    // Fall back to RemoteAddr (format: "IP:port")
    if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
        return host
    }
    
    return r.RemoteAddr
}


func IsHopByHopHeader(key string) bool {
    hopByHopHeaders := map[string]bool{
        "connection":          true,
        "keep-alive":          true,
        "proxy-authenticate":  true,
        "proxy-authorization": true,
        "te":                 true,
        "trailers":           true,
        "transfer-encoding":  true,
        "upgrade":            true,
    }
    
    return hopByHopHeaders[strings.ToLower(key)]
}