package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	gxc "gxclient-go"
	"gxclient-go/types"
	"testing"
)

const (
	testNetHttp = "https://testnet.gxchain.org"
	testNetWss  = "wss://testnet.gxchain.org"
	testFaucet  = "https://testnet.faucet.gxchain.org/account/register"

	testAccountName = "cli-wallet-test"
	testPri         = "5JsvYffKR8n4yNfCk36KkKFCzg6vo5fdBqqDJLavSifXSV9NABo"
	testPub         = "GXC58owosbFrudGVp8VCuMvDWpenx7AZSLwxEtAVqjWeqZ4YVLLWb"
)

func TestClient_Transfer(t *testing.T) {
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)

	result, err := client.Transfer("nathan", "ceshi 一下", "4.01 GXC", "GXC", true)
	require.NoError(t, err)
	str, _ := json.Marshal(*result)
	fmt.Println(string(str))
}

func TestClient_StakingCreate(t *testing.T) {
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetWss)
	require.Nil(t, err)

	result, err := client.CreateStaking("init0", 10.1, "1", "GXC", true)
	require.NoError(t, err)
	str, _ := json.Marshal(*result)
	fmt.Println(string(str))
}

func TestClient_StakingUpdate(t *testing.T) {
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)
	//owner
	owner, err := client.Database.GetAccount(testAccountName)

	to := "init0"
	toAccount, _ := client.Database.GetAccount(to)
	witness, _ := client.Database.GetWitnessByAccount(toAccount.ID.String())

	stakingObjects, err := client.Database.GetStakingObjects(owner.ID.String())
	var stakingObject types.StakingObject
	for _, ob := range stakingObjects {
		if ob.IsValid == true && ob.TrustNode.String() != witness.Id {
			stakingObject = *ob
			break
		}
	}

	result, err := client.UpdateStaking(to, stakingObject.ID.String(), "GXC", true)
	require.NoError(t, err)
	str, _ := json.Marshal(*result)
	fmt.Println(string(str))
}

func TestClient_StakingClaim(t *testing.T) {
	client, err := gxc.NewClient(testPri, testPri, testAccountName, testNetHttp)
	require.Nil(t, err)
	owner, err := client.Database.GetAccount(testAccountName)

	stakingObjects, err := client.Database.GetStakingObjects(owner.ID.String())
	var stakingObject types.StakingObject
	for _, ob := range stakingObjects {
		if ob.IsValid == false {
			stakingObject = *ob
		}
	}
	result, err := client.ClaimStaking(stakingObject.ID.String(), "GXC", true)
	require.NoError(t, err)
	str, _ := json.Marshal(*result)
	fmt.Println(string(str))
}
