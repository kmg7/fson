package netutils

import (
	"net"
	"net/http"
)

type AvailableIp struct {
	Name   string
	Ip     net.IP
	IsIPv4 bool
}

type DownloadFileHeaders struct {
	Name      string
	Size      int
	Resumable bool
}

func AvailableInterfaces() ([]AvailableIp, error) {
	ips := []AvailableIp{}
	infcs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, inf := range infcs {
		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}
		for _, v := range addrs {
			ip, _, err := net.ParseCIDR(v.String())
			if err != nil {
				continue
			}
			ips = append(ips, AvailableIp{Name: inf.Name, Ip: ip, IsIPv4: IsIPv4(ip)})
		}
	}
	return ips, nil

}

func DownloadHandler(
	w http.ResponseWriter,
	r *http.Request,
	p string,
	h map[string]string,
	isExistf func(string) (bool, error)) (bool, error) {
	ex, err := isExistf(p)
	if err != nil {
		http.NotFound(w, r)
		return false, err
	}
	if !ex {
		http.NotFound(w, r)
		return false, nil
	}
	SetHeaders(w, h)
	http.ServeFile(w, r, p)
	return true, nil

}

func IsIPv4(ip net.IP) bool {
	return ip.To4() != nil
}

func SetHeaders(w http.ResponseWriter, h map[string]string) {
	for k, v := range h {
		w.Header().Set(k, v)
	}
}

func ExtractHeaders(path string) (http.Header, error) {
	return nil, nil
}
