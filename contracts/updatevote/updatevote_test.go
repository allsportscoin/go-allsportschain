package updatevote

import (
	//"math/big"
	"testing"

	"github.com/allsportschain/go-allsportschain/accounts/abi/bind"
	//"github.com/allsportschain/go-allsportschain/accounts/abi/bind/backends"
	//"github.com/allsportschain/go-allsportschain/contracts/ens/contract"
	//"github.com/allsportschain/go-allsportschain/core"
	//"github.com/allsportschain/go-allsportschain/crypto"
	//"github.com/allsportschain/go-allsportschain/contracts/updatevote"
	//"fmt"
	"github.com/allsportschain/go-allsportschain/rpc"
	"github.com/allsportschain/go-allsportschain/socclient"
	"github.com/allsportschain/go-allsportschain/common"
	"context"
	//"github.com/allsportschain/go-allsportschain/node"
	//"fmt"
)

func TestUPDATEVOTE(t *testing.T) {
	//endpoint := "~/soft/socer/01/gsoc.ipc"
	endpoint := "/Users/bianhaoyi/soft/socer/01/gsoc.ipc"
	client, err := rpc.Dial(endpoint)
	if err != nil {
		t.Fatalf("error connecting to UPDATEVOTE API: %v", err)
	}
	socClient := socclient.NewClient(client)

	ct := new(context.Context)
	callOpts := bind.CallOpts{
		From: common.HexToAddress(""),
		Pending:false,
		Context:*ct,
	}
	udv ,err := NewUpdateVoteCaller(&callOpts, 1, socClient)
	if err != nil {
		t.Fatalf("can't deploy root registry: %v", err)
	}
	ret ,err := udv.GetVersion()
	//contractBackend.Commit()
	if err != nil {
		t.Fatalf("Test1 faild  %v", err)
	}

	print(ret,"\r\n")
}
