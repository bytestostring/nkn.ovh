package nknovh_engine

import (
		 "math/big"
		 "net"
		 "net/http"
		 "strings"
		 "time"
		 "math/rand"
		cr "crypto/rand"
		 "fmt"
		 "crypto/sha256"
		 "os"
		 "errors"
		)	

// IPv4 to int
func IP4toInt(ip string) int {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(net.ParseIP(ip).To4())
	return int(IPv4Int.Uint64())
}

func sleeping(iteration_time time.Duration, interval time.Duration, even bool) bool {
	var sleep time.Duration
	var t_even time.Duration = 120*time.Second
	var t_normal time.Duration = 60*time.Second
	var c_offset time.Duration = time.Duration(time.Now().Second())*time.Second + time.Duration(time.Now().Nanosecond())*time.Nanosecond

	if even == true {
		if e := time.Now().Minute() % 2; e == 0 {
			sleep = t_even - iteration_time
			if sleep <= 0 {
				time.Sleep(t_even-c_offset)
				return false
			}
			time.Sleep(sleep)
			return true
		}
	}
	sleep = interval*time.Second - iteration_time
	if sleep <= 0 {
		time.Sleep(t_normal-c_offset)
		return false
	}
	time.Sleep(sleep)
	return true
}

func RandBytes(n int) []byte {
	rand.Seed(time.Now().UnixNano())
	letterBytes := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$")
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return b
}


func GenRandomSHA256() (error, string) {
	h := sha256.New()
	b := make([]byte, 4096)
	_, err := cr.Read(b)
	if err != nil {
		return err, ""
	}
	h.Write(b)
	s := fmt.Sprintf("%x", h.Sum(nil))
	return nil, s
}

func getEtag(fn string) (error, int64) {
	f, err := os.Open(fn)
	if err != nil {
		return err, -1
	}
	defer f.Close()
	finfo, err := f.Stat()
	if err != nil {
		return err, -1
	}
	return nil, finfo.ModTime().Unix()
}

func FindStringInSlice(s []string, find string) int {
	for i, v := range s {
		if find == v {
			return i
		}
	}
	return len(s)
}

func FindIntInSlice(s []int, find int) int {
	for i, v := range s {
		if find == v {
			return i
		}
	}
	return len(s)
}

func IsPrivateIPv4(ip net.IP) (error, bool) {
	if ip4 := ip.To4(); ip4 != nil {
		return nil, ip4[0] == 10 ||
					ip4[0] == 127 ||
					(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
					(ip4[0] == 192 && ip4[1] == 168) ||
					(ip4[0] == 169 && ip4[1] == 254)
	}
	return errors.New("IP address is not IPv4"), false
}


func (o *NKNOVH) getIp(trusted []string, r *http.Request) (ip string, err error) {
	ip = ""
	forw := r.Header.Get("x-forwarded-for")
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	i := FindStringInSlice(trusted, host)
	if err != nil {
		return
	}
	if forw != "" {
		ips := strings.Split(forw, ",")
		if len(trusted) > i {
			ip = ips[0]
			if x := net.ParseIP(ip); x == nil {
				ip = host
			}
			return
		}
		ip = host
		return
	}
	realh := r.Header.Get("x-real-ip")
	if realh != "" {
		ips := strings.Split(realh, ",")
		if len(trusted) > i {
			ip = ips[0]
			if x := net.ParseIP(ip); x == nil {
				ip = host
			}
			return
		}
	}
	ip = host
	return
}
