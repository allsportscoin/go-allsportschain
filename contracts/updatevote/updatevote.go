package updatevote

import (
	//"strings"
	"github.com/allsportschain/go-allsportschain/accounts/abi/bind"
	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/contracts/updatevote/contract"
	//"github.com/allsportschain/go-allsportschain/core/types"
	//"github.com/allsportschain/go-allsportschain/crypto"
	"errors"
	"fmt"
)

var (
	MainNetAddress = common.HexToAddress("0x562a5606d9dd8415be0f4e05984b2c4e44b570e6")
	TestNetAddress = common.HexToAddress("0x86f767a9620c785bbc2bf2e530847e7a05273130")
	DevelopNetAddress = common.HexToAddress("0xf87a655d3e6da933211ade937da85738f82c8ba7")
)

// swarm domain name registry and resolver
type UPDATEVOTECaller struct {
	*contract.UPDATEVOTESessionCaller
	contractBackend bind.ContractBackend
}

// NewENS creates a struct exposing convenient high-level operations for interacting with
// the Ethereum Name Service.
func NewUpdateVoteCaller(callOpts *bind.CallOpts, netType uint64, contractBackend bind.ContractBackend) (*UPDATEVOTECaller, error) {
	var contractAddr common.Address
	if (netType == 1) {
		contractAddr = MainNetAddress
	} else if(netType == 3) {
		contractAddr = TestNetAddress
	} else if (netType == 2019) {
		contractAddr = DevelopNetAddress
	} else {
		err := errors.New(fmt.Sprintf("net Type invaild: %v", netType))
		return nil, err
	}


	updatevote, err := contract.NewUPDATEVOTECaller(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &UPDATEVOTECaller{
		&contract.UPDATEVOTESessionCaller{
			Contract:	updatevote,
			CallOpts: 	*callOpts,
		},
		contractBackend,
	}, nil
}

func DeployUPDATEVOTECaller(transactOpts *bind.TransactOpts, contractBackend bind.ContractBackend) (common.Address, *UPDATEVOTECaller, error) {
	addr, _, _, err := contract.DeployUPDATEVOTE(transactOpts, contractBackend)
	if err != nil {
		return addr, nil, err
	}
	callOpts := bind.CallOpts{
		From:transactOpts.From,
		Pending:false,
		Context:transactOpts.Context,
	}

	udv, err := contract.NewUPDATEVOTECaller(addr, contractBackend)
	if err != nil {
		return addr, nil, err
	}

	uudv := UPDATEVOTECaller{
		&contract.UPDATEVOTESessionCaller{
			Contract: udv,
			CallOpts: callOpts,
		},
		contractBackend,
	}

	return addr, &uudv, nil
}

func (self *UPDATEVOTECaller) GetVersion() (string, error) {
	ret,err := self.Contract.GetVersion(&self.CallOpts)
	if (err != nil) {
		return "", err
	}
	return ret, nil
}


