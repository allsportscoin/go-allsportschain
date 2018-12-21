package sns


//go:generate abigen --sol contract/SOCResolver.sol --exc contract/AbstractSOCResolver.sol:AbstractSOCResolver --pkg contract --out contract/socresolver.go
//go:generate abigen --sol contract/SNS.sol --exc contract/AbstractSNS.sol:AbstractSNS --pkg contract --out contract/sns.go

import (
	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/contracts/sns/contract"
	"github.com/allsportschain/go-allsportschain/accounts/abi/bind"
)

var (
	SuffixName = string(".soc")

	MainNetSNSAddress = common.HexToAddress("0xb2756edb776b8223dacb95eaceea96c29a7e79f4")
	MainNetSOCAddress = common.HexToAddress("0xb2756edb776b8223dacb95eaceea96c29a7e79f4")

	TestNetSNSAddress = common.HexToAddress("0x1263bddc0f9b45b5b60da643a3aeb92c07ef99cd")
	TestNetSOCAddress = common.HexToAddress("0xaf3c93eeed2415e83544f2df07e9ee9027c71810")

	//0x788a402f20e3fa4d0a3d03a9c9dc159f6d60755adbab37a99b48c00baba47ae9
	DevelopNetSNSAddress = common.HexToAddress("0xaf9b04877a5e1a0b5426bd7db2ef141203aca3bb")
	//0x42f04769011ad0ea27da910c9ddef5246f39e510e4c7cc823adfda0e1488ba66
	DevelopNetSOCAddress = common.HexToAddress("0x6952c8eaf85a5ebfd050e02054a964126083251f")

	DefaultAddress = common.Address{}
)

type SNSCaller struct {
	//*contract.SNS
	contract.SNSCallerSession
}

// NewENS creates a struct exposing convenient high-level operations for interacting with
// the Ethereum Name Service.
func NewSNSCaller(callOpts *bind.CallOpts, contractBackend bind.ContractBackend, networkId uint64) (*SNSCaller, error) {
	snsContractAddr := common.Address{}
	switch networkId {
	case 1:
		snsContractAddr = MainNetSNSAddress
	case 3:
		snsContractAddr = TestNetSNSAddress
	case 2019:
		snsContractAddr = DevelopNetSNSAddress
	default:
		snsContractAddr = DefaultAddress
	}
	//if (netType == 1) {
	//	contractAddr = MainNetAddress
	//} else if(netType == 3) {
	//	contractAddr = TestNetAddress
	//} else if (netType == 2019) {
	//	contractAddr = DevelopNetAddress
	//} else {
	//	err := errors.New(fmt.Sprintf("net Type invaild: %v", netType))
	//	return nil, err
	//}

	sns, err := contract.NewSNSCaller(snsContractAddr, contractBackend)
	if (err != nil) {
		return nil, err
	}

	return &SNSCaller{
		contract.SNSCallerSession{
			Contract: sns,
			CallOpts: *callOpts,
		},
	}, nil

}

func (self *SNSCaller) GetAddrByAlias(aliasName string) (common.Address, error) {
	ret, err := self.Contract.GetAddrByAliasAndSuffix(&self.CallOpts, SuffixName, aliasName)
	if (err != nil) {
		//log.Info("alias not exist,by alias.",aliasName)
		return common.Address{}, err
	}
	return ret, nil
}

func (self *SNSCaller)GetAliasByAddr(addr common.Address) (string, error) {
	ret, err := self.Contract.GetAliasByAddrAndSuffix(&self.CallOpts, SuffixName, addr)
	if (err != nil) {
		//log.Info("alias not exist by addr.",addr.Hex())
		return "", err
	}
	return ret,nil
}