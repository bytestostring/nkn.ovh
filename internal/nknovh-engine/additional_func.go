package nknovh_engine

import (
		 "math/big"
		 "net"
		 "time"
		 "math/rand"
		)	

// IPv4 to int
func IP4toInt(IPv4Addr string) int {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(net.ParseIP(IPv4Addr).To4())
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