package main

import (
	"flag"
	"fmt"
	"gxclient-go/offline"
)

func main() {
	activePriWif := flag.String("key", "", "wif of active private key")
	unSignedHex := flag.String("txHex", "", "hex of unSigned transaction")
	chainId := flag.String("chainId", "", "chainId")
	flag.Parse()

	sig, err := offline.OfflineSign(*activePriWif, *unSignedHex, *chainId)
	if err != nil {
		fmt.Println("sign failed :", err)
	}
	fmt.Println(sig)
}
