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

func sleeping(iteration_time time.Duration, interval time.Duration) {
	sleep := interval*time.Second - iteration_time
	time.Sleep(sleep)
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