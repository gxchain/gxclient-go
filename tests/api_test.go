package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	gxc "gxclient-go"
	"testing"
)

func TestApi_GetChainId(t *testing.T) {
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)
	chainId, err := client.Database.GetChainId()
	fmt.Println(chainId)
}

// GetDynamicGlobalProperties Gets dynamic global properties of current blockchain
func TestApi_GetDynamicGlobalProperties(t *testing.T) {
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)

	properties, err := client.Database.GetDynamicGlobalProperties()
	require.Nil(t, err)
	s, _ := json.Marshal(properties)
	fmt.Println(string(s))
}

// GetBlock return a block by the given block number
func TestApi_GetBlock(t *testing.T) {
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)
	block, err := client.Database.GetBlock(22039351)
	str, _ := json.Marshal(*block)
	fmt.Println(string(str))
}

func TestApi_GetObjects(t *testing.T) {
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)

	o1 := "1.3.1"
	o2 := "1.3.11"
	obs, err := client.Database.GetObjects(o1, o2)
	for _, o := range obs {
		str := string([]byte(o)[:])
		fmt.Println(str)
	}

	ob, err := client.Database.GetObject(o1)
	str := string([]byte(ob)[:])
	fmt.Println(str)
}

func TestApi_GetAccounts(t *testing.T) {
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)

	a1 := "null-account"
	a2 := "init0"
	accs, err := client.Database.GetAccounts(a1, a2)
	require.Nil(t, err)
	for _, acc := range accs {
		str, _ := json.Marshal(*acc)
		fmt.Println(string(str))
	}

	acc, err := client.Database.GetAccount("t")
	require.Nil(t, err)
	str, _ := json.Marshal(*acc)
	fmt.Println(string(str))
}

func TestApi_GetAccountBalance(t *testing.T) {
	accountName := "cli-wallet-test-1"
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)
	databaseApi := client.Database

	account, err := client.Database.GetAccount(accountName)
	require.Nil(t, err)

	assets, err := client.Database.GetAsset("BDB")
	require.Nil(t, err)
	accountBalances, err := databaseApi.GetAccountBalances(account.ID.String(), assets.ID.String())
	require.Nil(t, err)
	for _, accBalance := range accountBalances {
		fmt.Printf("assetsId: %s ,amount: %d \n", accBalance.AssetID, accBalance.Amount)
	}
}

func TestApi_GetAccountBalances(t *testing.T) {
	accountName := "dev"
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)
	databaseApi := client.Database

	account, err := client.Database.GetAccount(accountName)
	require.Nil(t, err)

	accountBalances, err := databaseApi.GetAccountBalances(account.ID.String())
	require.Nil(t, err)

	for _, accBalance := range accountBalances {
		fmt.Printf("assetsId: %s ,amount: %d \n", accBalance.AssetID, accBalance.Amount)
	}
}

func TestApi_GetAccountsByPublicKeys(t *testing.T) {
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)

	pub1 := "GXC6dwwmF98DDrRj3R6ZqvYQanTUcEy6QGrqtbhwpDRGtrd6P9sob" //dev
	pub2 := "GXC8K2LPx4WFv2E4twhjFaom7AjDSK6dPj2N3z42kfQPbcb2aeVUH"
	pub3 := "GXC58owosbFrudGVp8VCuMvDWpenx7AZSLwxEtAVqjWeqZ4YVLLWb"
	accs, err := client.Database.GetAccountsByPublicKeys(pub1, pub2, pub3)
	require.Nil(t, err)

	for _, acc := range *accs {
		for _, a := range acc {
			fmt.Println(a)
		}
	}

	a, err := client.Database.GetAccountsByPublicKey(pub3)
	fmt.Println(a)

}

func TestApi_GetAssets(t *testing.T) {
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)
	assets, err := client.Database.GetAssets("GXC", "BDB", "NULL")
	asset, err := client.Database.GetAsset("1.3.100")
	fmt.Println(assets)
	fmt.Println(asset)
}

func TestApi_Simple(t *testing.T) {
	//amount, _ := strconv.ParseFloat("4.1", 64)
	//a := decimal.NewFromFloat(amount).Mul(decimal.NewFromFloat(math.Pow10(int(uint8(5)))))
	//b := a.IntPart()
	//println(b)

	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)

	tx, err := client.Database.GeTransactionExtByTxid("0101813c34fb033b7ba7a30c675bfa1b949357d8")
	str1, _ := json.Marshal(tx)
	fmt.Println(string(str1))

	//tx, err := client.Database.GetTransactionByTxid("647524E8668D3484A764009B0AC90394CF64ED3C")
	//str1, _ := json.Marshal(tx)
	//fmt.Println(string(str1))
	//
	//s, err := client.Database.GetStakingPrograms()
	//fmt.Println(s)
	//
	////GetBlockHeader
	//header, _ := client.Database.GetBlockHeader(22039352)
	//str, _ := json.Marshal(*header)
	//fmt.Println(string(str))

}
