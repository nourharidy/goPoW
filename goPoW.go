package main

import "github.com/majestrate/cryptonight"
import "fmt"
import "os"
import "strings"
import "strconv"
import "runtime"
import "sync"
import "github.com/logrusorgru/aurora"
import "crypto/sha256"

func main() {
	nonce := make(chan int)

	var wg sync.WaitGroup
	wg.Add(1)

	if len(os.Args) < 3 {
		fmt.Printf("You must enter a hash as a first argument and the amount of work as a 2nd argument. \n")
	} else {
		cpus := runtime.NumCPU();
		if len(os.Args) > 3 {
			newCPUs, e := strconv.Atoi(os.Args[3]);
			cpus = newCPUs
			if e != nil {
				fmt.Println(e)
			}
		}
		algo := "cryptonight"
		if len(os.Args) > 4 {
			algo = os.Args[4]
		}
		for i := 0; i <= cpus; i++ {
			go pow(nonce, &wg, algo, os.Args[1], os.Args[2])
		}
		nonce <- 0
		runtime.GOMAXPROCS(cpus)
		fmt.Printf("Using %d CPUs\n", cpus)
		wg.Wait()
	}
}

func pow(nonce chan int, wg *sync.WaitGroup, algo string, text string, difficulty string) {
	var target string;
	for len(target) == 0 {
		currentNonce := <- nonce
		currentNonce++
		nonce <- currentNonce
		input := strings.Join([]string{text, string(currentNonce)}, "")
		var hash [32]byte;
		if algo == "cryptonight" {
			hash = cryptonight.HashBytes([]byte(input))
		}else if algo == "sha256" {
			hash = sha256.Sum256([]byte(input))
		}else{
			// Default algo: Cryptonight
			hash = cryptonight.HashBytes([]byte(input))
		}
		fmt.Printf("Nonce %d: %x\n", currentNonce, hash)
		zeros, e := strconv.Atoi(difficulty)
		if e != nil {
			fmt.Println(e)
		}
		if fmt.Sprintf("%x", hash)[0:zeros] == strings.Repeat("0", zeros ) {
			target = fmt.Sprintf("%x", hash)
			fmt.Printf("\nFound match!\n\nNonce: %d\nHash: %x\n", aurora.Blue(currentNonce), aurora.Green(hash))
			wg.Done()
		}
	}
}