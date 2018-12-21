package contract

import (
	"strings"

	//ethereum "github.com/allsportschain/go-allsportschain"
	"github.com/allsportschain/go-allsportschain/accounts/abi"
	"github.com/allsportschain/go-allsportschain/accounts/abi/bind"
	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/core/types"
	//"github.com/allsportschain/go-allsportschain/event"
)

// UPDATEVOTEABI is the input ABI used to generate the binding from.
//const UPDATEVOTEABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"resolver\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"bytes32\"},{\"name\":\"label\",\"type\":\"bytes32\"},{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"setSubnodeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"bytes32\"},{\"name\":\"ttl\",\"type\":\"uint64\"}],\"name\":\"setTTL\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"ttl\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"bytes32\"},{\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"setResolver\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"bytes32\"},{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"label\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"NewOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"NewResolver\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"ttl\",\"type\":\"uint64\"}],\"name\":\"NewTTL\",\"type\":\"event\"}]"

const UPDATEVOTEABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getProposalState\",\"outputs\":[{\"name\":\"ifinshNode\",\"type\":\"uint256\"},{\"name\":\"idesc\",\"type\":\"string\"},{\"name\":\"iversion\",\"type\":\"bytes32\"},{\"name\":\"needPassNum\",\"type\":\"uint256\"},{\"name\":\"icondidate\",\"type\":\"address[]\"},{\"name\":\"pass\",\"type\":\"bool[]\"},{\"name\":\"len\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"name\":\"version\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"pass\",\"type\":\"bool\"}],\"name\":\"voteUpdate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"ifinshNode\",\"type\":\"uint256\"},{\"name\":\"version\",\"type\":\"bytes32\"},{\"name\":\"icandidate\",\"type\":\"address[]\"},{\"name\":\"candidateLen\",\"type\":\"uint256\"},{\"name\":\"passNum\",\"type\":\"uint256\"},{\"name\":\"idesc\",\"type\":\"string\"}],\"name\":\"putProposal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"cancalPoroposal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"}]"


// UPDATEVOTEBin is the compiled bytecode used for deploying new contracts.
const UPDATEVOTEBin = `0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060018060006101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550600160009054906101000a900467ffffffffffffffff16600060146101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060006002819055506000600560006101000a81548160ff0219169083151502179055506000600381905550610e7a806100f36000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630d8e6e2c1461006f57806364b13450146100ae578063aade375b146100dd578063bcdcc9ef1461022d578063d1292f8a1461030b575b005b34801561007b57600080fd5b50610084610322565b604051808267ffffffffffffffff1667ffffffffffffffff16815260200191505060405180910390f35b3480156100ba57600080fd5b506100db6004803603810190808035151590602001909291905050506103cd565b005b3480156100e957600080fd5b506100f261069d565b60405180888152602001806020018767ffffffffffffffff1667ffffffffffffffff168152602001868152602001806020018060200185815260200184810384528a818151815260200191508051906020019080838360005b8381101561016657808201518184015260208101905061014b565b50505050905090810190601f1680156101935780820380516001836020036101000a031916815260200191505b50848103835287818151815260200191508051906020019060200280838360005b838110156101cf5780820151818401526020810190506101b4565b50505050905001848103825286818151815260200191508051906020019060200280838360005b838110156102115780820151818401526020810190506101f6565b505050509050019a505050505050505050505060405180910390f35b34801561023957600080fd5b5061030960048036038101908080359060200190929190803567ffffffffffffffff169060200190929190803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192908035906020019092919080359060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610939565b005b34801561031757600080fd5b50610320610c87565b005b60008060149054906101000a900467ffffffffffffffff1667ffffffffffffffff16600160009054906101000a900467ffffffffffffffff1667ffffffffffffffff161415801561037557506002544310155b801561038d5750600560009054906101000a900460ff165b156103b057600060149054906101000a900467ffffffffffffffff1690506103ca565b600160009054906101000a900467ffffffffffffffff1690505b90565b6000806000600160009054906101000a900467ffffffffffffffff1667ffffffffffffffff16600060149054906101000a900467ffffffffffffffff1667ffffffffffffffff1614158015610423575043600254115b151561042e57600080fd5b6000925060009150600090505b600680549050811015610644573373ffffffffffffffffffffffffffffffffffffffff1660068281548110151561046e57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141561058e57831561052457600060149054906101000a900467ffffffffffffffff1667ffffffffffffffff16600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610589565b600160009054906101000a900467ffffffffffffffff1667ffffffffffffffff16600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505b600192505b600060149054906101000a900467ffffffffffffffff1667ffffffffffffffff16600760006006848154811015156105c257fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414156106375781806001019250505b808060010191505061043b565b82151561065057600080fd5b6003548210151561067b576001600560006101000a81548160ff021916908315150217905550610697565b6000600560006101000a81548160ff0219169083151502179055505b50505050565b6000606060008060608060008060149054906101000a900467ffffffffffffffff1667ffffffffffffffff16600160009054906101000a900467ffffffffffffffff1667ffffffffffffffff1614156106f95760009650610930565b600680548060200260200160405190810160405280929190818152602001828054801561077b57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610731575b505050505092506006805490506040519080825280602002602001820160405280156107b65781602001602082028038833980820191505090505b509150600090505b825181101561086f57600060149054906101000a900467ffffffffffffffff1667ffffffffffffffff166007600085848151811015156107fa57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414828281518110151561084e57fe5b906020019060200201901515908115158152505080806001019150506107be565b60025496506003549350600060149054906101000a900467ffffffffffffffff16945060048054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156109285780601f106108fd57610100808354040283529160200191610928565b820191906000526020600020905b81548152906001019060200180831161090b57829003601f168201915b505050505095505b90919293949596565b60003373ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614151561099657600080fd5b43600254101515610a8d57600560009054906101000a900460ff1615610a1557600060149054906101000a900467ffffffffffffffff16600160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506000600560006101000a81548160ff021916908315150217905550610a55565b600160009054906101000a900467ffffffffffffffff16600060146101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055505b60006002819055506000600381905550602060405190810160405280600081525060049080519060200190610a8b929190610d7d565b505b600060149054906101000a900467ffffffffffffffff1667ffffffffffffffff16600160009054906101000a900467ffffffffffffffff1667ffffffffffffffff16148015610b065750600160009054906101000a900467ffffffffffffffff1667ffffffffffffffff168667ffffffffffffffff1614155b8015610b1157504387115b1515610b85576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600b8152602001807f6261642056657273696f6e00000000000000000000000000000000000000000081525060200191505060405180910390fd5b8660028190555085600060146101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506000600681610bc49190610dfd565b5060009050600090505b83811015610c605760068582815181101515610be657fe5b9060200190602002015190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550508080600101915050610bce565b826003819055508160049080519060200190610c7d929190610d7d565b5050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141515610ce257600080fd5b43600254111515610cf257600080fd5b600160009054906101000a900467ffffffffffffffff16600060146101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506000600560006101000a81548160ff0219169083151502179055506000600281905550602060405190810160405280600081525060049080519060200190610d7a929190610d7d565b50565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610dbe57805160ff1916838001178555610dec565b82800160010185558215610dec579182015b82811115610deb578251825591602001919060010190610dd0565b5b509050610df99190610e29565b5090565b815481835581811115610e2457818360005260206000209182019101610e239190610e29565b5b505050565b610e4b91905b80821115610e47576000816000905550600101610e2f565b5090565b905600a165627a7a723058203a261d28533bd2330d0c38307d699e7d62c9e2b3b2fabe7c1b28f5b74431a1b20029`

// UPDATEVOTE is an auto generated Go binding around an Ethereum contract.
type UPDATEVOTE struct {
	UPDATEVOTECaller     // Read-only binding to the contract
	UPDATEVOTETransactor // Write-only binding to the contract
	UPDATEVOTEFilterer   // Log filterer for contract events
}

// UPDATEVOTECaller is an auto generated read-only Go binding around an Ethereum contract.
type UPDATEVOTECaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UPDATEVOTETransactor is an auto generated write-only Go binding around an Ethereum contract.
type UPDATEVOTETransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UPDATEVOTEFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UPDATEVOTEFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NewUPDATEVOTE creates a new instance of UPDATEVOTE, bound to a specific deployed contract.
func NewUPDATEVOTE(address common.Address, backend bind.ContractBackend) (*UPDATEVOTE, error) {
	contract, err := bindUPDATEVOTE(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UPDATEVOTE{UPDATEVOTECaller: UPDATEVOTECaller{contract: contract}, UPDATEVOTETransactor: UPDATEVOTETransactor{contract: contract}, UPDATEVOTEFilterer: UPDATEVOTEFilterer{contract: contract}}, nil
}

// NewUPDATEVOTECaller creates a new read-only instance of UPDATEVOTE, bound to a specific deployed contract.
func NewUPDATEVOTECaller(address common.Address, caller bind.ContractCaller) (*UPDATEVOTECaller, error) {
	contract, err := bindUPDATEVOTE(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UPDATEVOTECaller{contract: contract}, nil
}

// bindUPDATEVOTE binds a generic wrapper to an already deployed contract.
func bindUPDATEVOTE(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UPDATEVOTEABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// DeployUPDATEVOTE deploys a new Ethereum contract, binding an instance of UPDATEVOTE to it.
func DeployUPDATEVOTE(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UPDATEVOTE, error) {
	parsed, err := abi.JSON(strings.NewReader(UPDATEVOTEABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(UPDATEVOTEBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UPDATEVOTE{UPDATEVOTECaller: UPDATEVOTECaller{contract: contract}, UPDATEVOTETransactor: UPDATEVOTETransactor{contract: contract}, UPDATEVOTEFilterer: UPDATEVOTEFilterer{contract: contract}}, nil
}

// ENSSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UPDATEVOTESessionCaller struct {
	Contract     *UPDATEVOTECaller       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
}

func (_UPD *UPDATEVOTECaller) GetVersion(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0;
	err := _UPD.contract.Call(opts, out, "getVersion")
	return *ret0, err
}


