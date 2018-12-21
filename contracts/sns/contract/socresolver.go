// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"strings"

	"github.com/allsportschain/go-allsportschain/accounts/abi"
	"github.com/allsportschain/go-allsportschain/accounts/abi/bind"
	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/core/types"
)

// SOCResolverABI is the input ABI used to generate the binding from.
const SOCResolverABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"withDraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getSuffix\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_aliasName\",\"type\":\"string\"}],\"name\":\"getAddrByAlias\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_aliasName\",\"type\":\"string\"},{\"name\":\"_aliasAddr\",\"type\":\"address\"}],\"name\":\"setAlias\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_aliasName\",\"type\":\"string\"},{\"name\":\"_aliasAddr\",\"type\":\"address\"}],\"name\":\"setVipAlias\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getAliasByAddr\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"Suffix\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportSNSInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_suffix\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// SOCResolverBin is the compiled bytecode used for deploying new contracts.
const SOCResolverBin = `0x60806040523480156200001157600080fd5b5060405162001a7238038062001a7283398101604052805160008054600160a060020a03191633600160a060020a03161790550180516200005a90600190602084019062000062565b505062000107565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620000a557805160ff1916838001178555620000d5565b82800160010185558215620000d5579182015b82811115620000d5578251825591602001919060010190620000b8565b50620000e3929150620000e7565b5090565b6200010491905b80821115620000e35760008155600101620000ee565b90565b61195b80620001176000396000f3006080604052600436106100a35763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630a67d2c781146100a85780631c0e1d09146100cb578063351af6a0146101555780633912718c146101ca57806345c268521461022157806351119ae1146102855780638da5cb5b146102a657806398e77135146102bb578063c67b08d2146102d0578063df1827df14610306575b600080fd5b3480156100b457600080fd5b506100c9600160a060020a036004351661036b565b005b3480156100d757600080fd5b506100e06103d9565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561011a578181015183820152602001610102565b50505050905090810190601f1680156101475780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561016157600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526101ae94369492936024939284019190819084018382808284375094975061046f9650505050505050565b60408051600160a060020a039092168252519081900360200190f35b6040805160206004803580820135601f81018490048402850184019095528484526100c994369492936024939284019190819084018382808284375094975050509235600160a060020a0316935061067292505050565b34801561022d57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100c994369492936024939284019190819084018382808284375094975050509235600160a060020a03169350610fb692505050565b34801561029157600080fd5b506100e0600160a060020a036004351661130f565b3480156102b257600080fd5b506101ae6114f3565b3480156102c757600080fd5b506100e0611502565b3480156102dc57600080fd5b506102f2600160e060020a03196004351661158f565b604080519115158252519081900360200190f35b34801561031257600080fd5b5061031b611694565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561035757818101518382015260200161033f565b505050509050019250505060405180910390f35b60005433600160a060020a0390811691161461038657600080fd5b600160a060020a038116151561039b57600080fd5b604051600160a060020a0380831691309091163180156108fc02916000818181858888f193505050501580156103d5573d6000803e3d6000fd5b5050565b60018054604080516020601f600260001961010087891615020190951694909404938401819004810282018101909252828152606093909290918301828280156104645780601f1061043957610100808354040283529160200191610464565b820191906000526020600020905b81548152906001019060200180831161044757829003601f168201915b505050505090505b90565b6000610479611881565b6003836040518082805190602001908083835b602083106104ab5780518252601f19909201916020918201910161048c565b51815160209384036101000a600019018019909216911617905292019485525060405193849003019092206001015460a060020a900460ff161515915061053e9050576040805160e560020a62461bcd02815260206004820152601560248201527f616c696173206e616d65206e6f74206578697374730000000000000000000000604482015290519081900360640190fd5b6003836040518082805190602001908083835b602083106105705780518252601f199092019160209182019101610551565b518151600019602094850361010090810a8201928316921993909316919091179092529490920196875260408051978890038201882080546080601f600260018416159099029096019091169690960493840183900490920288018501905260608701828152909550869450928592508401828280156106315780601f1061060657610100808354040283529160200191610631565b820191906000526020600020905b81548152906001019060200180831161061457829003601f168201915b505050918352505060019190910154600160a060020a03811660208084019190915260a060020a90910460ff16151560409092019190915201519392505050565b60018054604080516020601f6002600019610100878916150201909516949094049384018190048102820181019092528281528593909290918301828280156106fc5780601f106106d1576101008083540402835291602001916106fc565b820191906000526020600020905b8154815290600101906020018083116106df57829003601f168201915b5050845184518694508593509091506000806006830184118015610722575082600d0184105b151561079e576040805160e560020a62461bcd02815260206004820152602660248201527f646f6d61696e206c656e6774682065786365707420737566666978206b65657060448201527f7320372d31320000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600091505b8282101561088057846001838503038151811015156107be57fe5b90602001015160f860020a900460f860020a02600160f860020a031916866001848703038151811015156107ee57fe5b60209101015160f860020a90819004027fff000000000000000000000000000000000000000000000000000000000000001614610875576040805160e560020a62461bcd02815260206004820152601460248201527f696e76616c696420737566666978206572726f72000000000000000000000000604482015290519081900360640190fd5b6001909101906107a3565b85517f410000000000000000000000000000000000000000000000000000000000000090879060009081106108b157fe5b90602001015160f860020a900460f860020a02600160f860020a03191610158015610928575085517f5a00000000000000000000000000000000000000000000000000000000000000908790600090811061090857fe5b90602001015160f860020a900460f860020a02600160f860020a03191611155b806109d6575085517f6100000000000000000000000000000000000000000000000000000000000000908790600090811061095f57fe5b90602001015160f860020a900460f860020a02600160f860020a031916101580156109d6575085517f7a0000000000000000000000000000000000000000000000000000000000000090879060009081106109b657fe5b90602001015160f860020a900460f860020a02600160f860020a03191611155b1515610a52576040805160e560020a62461bcd02815260206004820152602b60248201527f496e697469616c20636861726163746572206d75737420626520656e676c697360448201527f6820616c7068616265742e000000000000000000000000000000000000000000606482015290519081900360840190fd5b5060015b828403811015610cbc5785517f300000000000000000000000000000000000000000000000000000000000000090879083908110610a9057fe5b90602001015160f860020a900460f860020a02600160f860020a03191610158015610b06575085517f390000000000000000000000000000000000000000000000000000000000000090879083908110610ae657fe5b90602001015160f860020a900460f860020a02600160f860020a03191611155b80610bb2575085517f410000000000000000000000000000000000000000000000000000000000000090879083908110610b3c57fe5b90602001015160f860020a900460f860020a02600160f860020a03191610158015610bb2575085517f5a0000000000000000000000000000000000000000000000000000000000000090879083908110610b9257fe5b90602001015160f860020a900460f860020a02600160f860020a03191611155b80610c5e575085517f610000000000000000000000000000000000000000000000000000000000000090879083908110610be857fe5b90602001015160f860020a900460f860020a02600160f860020a03191610158015610c5e575085517f7a0000000000000000000000000000000000000000000000000000000000000090879083908110610c3e57fe5b90602001015160f860020a900460f860020a02600160f860020a03191611155b1515610cb4576040805160e560020a62461bcd02815260206004820152601860248201527f696e76616c696420636861726163746572206572726f722e0000000000000000604482015290519081900360640190fd5b600101610a56565b60038a6040518082805190602001908083835b60208310610cee5780518252601f199092019160209182019101610ccf565b51815160209384036101000a600019018019909216911617905292019485525060405193849003019092206001015460a060020a900460ff16159150610d809050576040805160e560020a62461bcd02815260206004820152601460248201527f616c69617320616c726561647920657869737473000000000000000000000000604482015290519081900360640190fd5b600160a060020a0389161515610d94573398505b600160a060020a0389166000908152600460205260409020546002600019610100600184161502019091160415610e3b576040805160e560020a62461bcd02815260206004820152602560248201527f6164647265737320616c72656164792062696e6420746f20616e6f746865722060448201527f616c696173000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b60028054600181018083556000929092528b51610e7f917f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace019060208e01906118a1565b5050600160a060020a03891660009081526004602090815260409091208b51610eaa928d01906118a1565b506060604051908101604052808b81526020018a600160a060020a031681526020016001151581525060038b6040518082805190602001908083835b60208310610f055780518252601f199092019160209182019101610ee6565b51815160209384036101000a6000190180199092169116179052920194855250604051938490038101909320845180519194610f46945085935001906118a1565b50602082015160019091018054604090930151151560a060020a0274ff000000000000000000000000000000000000000019600160a060020a0390931673ffffffffffffffffffffffffffffffffffffffff19909416939093179190911691909117905550505050505050505050565b60005433600160a060020a03908116911614610fd157600080fd5b6003826040518082805190602001908083835b602083106110035780518252601f199092019160209182019101610fe4565b51815160209384036101000a600019018019909216911617905292019485525060405193849003019092206001015460a060020a900460ff161591506110959050576040805160e560020a62461bcd02815260206004820152601460248201527f616c69617320616c726561647920657869737473000000000000000000000000604482015290519081900360640190fd5b600160a060020a03811615156110f5576040805160e560020a62461bcd02815260206004820152601960248201527f656d7074792061646472657373206e6f7420616c6c6f77656400000000000000604482015290519081900360640190fd5b600160a060020a038116600090815260046020526040902054600260001961010060018416150201909116041561119c576040805160e560020a62461bcd02815260206004820152602560248201527f6164647265737320616c72656164792062696e6420746f20616e6f746865722060448201527f616c696173000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600280546001810180835560009290925283516111e0917f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace019060208601906118a1565b5050600160a060020a0381166000908152600460209081526040909120835161120b928501906118a1565b5060606040519081016040528083815260200182600160a060020a03168152602001600115158152506003836040518082805190602001908083835b602083106112665780518252601f199092019160209182019101611247565b51815160209384036101000a60001901801990921691161790529201948552506040519384900381019093208451805191946112a7945085935001906118a1565b50602082015160019091018054604090930151151560a060020a0274ff000000000000000000000000000000000000000019600160a060020a0390931673ffffffffffffffffffffffffffffffffffffffff1990941693909317919091169190911790555050565b6060611319611881565b600160a060020a038316600090815260046020526040812054600260001961010060018416150201909116041161139a576040805160e560020a62461bcd02815260206004820152601960248201527f616c696173206f66205f61646472206e6f742065786973747300000000000000604482015290519081900360640190fd5b60036004600085600160a060020a0316600160a060020a03168152602001908152602001600020604051808280546001816001161561010002031660029004801561141c5780601f106113fa57610100808354040283529182019161141c565b820191906000526020600020905b815481529060010190602001808311611408575b505092835250506040805160209281900383018120805460026001821615610100026000190190911604601f810185900490940282016080908101909352606082018481529193909284929184918401828280156114bb5780601f10611490576101008083540402835291602001916114bb565b820191906000526020600020905b81548152906001019060200180831161149e57829003601f168201915b505050918352505060019190910154600160a060020a038116602083015260a060020a900460ff161515604090910152519392505050565b600054600160a060020a031681565b60018054604080516020600284861615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156115875780601f1061155c57610100808354040283529160200191611587565b820191906000526020600020905b81548152906001019060200180831161156a57829003601f168201915b505050505081565b60007f351af6a000000000000000000000000000000000000000000000000000000000600160e060020a0319831614806115f257507f51119ae100000000000000000000000000000000000000000000000000000000600160e060020a03198316145b8061162657507f3912718c00000000000000000000000000000000000000000000000000000000600160e060020a03198316145b8061165a57507fc67b08d200000000000000000000000000000000000000000000000000000000600160e060020a03198316145b8061168e57507f1c0e1d0900000000000000000000000000000000000000000000000000000000600160e060020a03198316145b92915050565b60408051600580825260c082019092526060918291906020820160a08038833901905050604080517f736574416c69617328737472696e672c6164647265737329000000000000000081529051908190036018019020815191925090829060009081106116fd57fe5b600160e060020a0319909216602092830290910190910152604080517f676574416464724279416c69617328737472696e6729000000000000000000008152905190819003601601902081518290600190811061175657fe5b600160e060020a0319909216602092830290910190910152604080517f676574416c696173427941646472286164647265737329000000000000000000815290519081900360170190208151829060029081106117af57fe5b600160e060020a0319909216602092830290910190910152604080517f737570706f7274534e53496e74657266616365286279746573342900000000008152905190819003601b01902081518290600390811061180857fe5b600160e060020a0319909216602092830290910190910152604080517f67657453756666697828290000000000000000000000000000000000000000008152905190819003600b01902081518290600490811061186157fe5b600160e060020a03199092166020928302909101909101529050805b5090565b604080516060818101835281526000602082018190529181019190915290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106118e257805160ff191683800117855561190f565b8280016001018555821561190f579182015b8281111561190f5782518255916020019190600101906118f4565b5061187d9261046c9250905b8082111561187d576000815560010161191b5600a165627a7a72305820e509ac07afafbd4a21636c316a0d6b942f33579f66f0d5064a11383e15c908bd0029`

// DeploySOCResolver deploys a new Ethereum contract, binding an instance of SOCResolver to it.
func DeploySOCResolver(auth *bind.TransactOpts, backend bind.ContractBackend, _suffix string) (common.Address, *types.Transaction, *SOCResolver, error) {
	parsed, err := abi.JSON(strings.NewReader(SOCResolverABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SOCResolverBin), backend, _suffix)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SOCResolver{SOCResolverCaller: SOCResolverCaller{contract: contract}, SOCResolverTransactor: SOCResolverTransactor{contract: contract}, SOCResolverFilterer: SOCResolverFilterer{contract: contract}}, nil
}

// SOCResolver is an auto generated Go binding around an Ethereum contract.
type SOCResolver struct {
	SOCResolverCaller     // Read-only binding to the contract
	SOCResolverTransactor // Write-only binding to the contract
	SOCResolverFilterer   // Log filterer for contract events
}

// SOCResolverCaller is an auto generated read-only Go binding around an Ethereum contract.
type SOCResolverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SOCResolverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SOCResolverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SOCResolverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SOCResolverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SOCResolverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SOCResolverSession struct {
	Contract     *SOCResolver      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SOCResolverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SOCResolverCallerSession struct {
	Contract *SOCResolverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// SOCResolverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SOCResolverTransactorSession struct {
	Contract     *SOCResolverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SOCResolverRaw is an auto generated low-level Go binding around an Ethereum contract.
type SOCResolverRaw struct {
	Contract *SOCResolver // Generic contract binding to access the raw methods on
}

// SOCResolverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SOCResolverCallerRaw struct {
	Contract *SOCResolverCaller // Generic read-only contract binding to access the raw methods on
}

// SOCResolverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SOCResolverTransactorRaw struct {
	Contract *SOCResolverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSOCResolver creates a new instance of SOCResolver, bound to a specific deployed contract.
func NewSOCResolver(address common.Address, backend bind.ContractBackend) (*SOCResolver, error) {
	contract, err := bindSOCResolver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SOCResolver{SOCResolverCaller: SOCResolverCaller{contract: contract}, SOCResolverTransactor: SOCResolverTransactor{contract: contract}, SOCResolverFilterer: SOCResolverFilterer{contract: contract}}, nil
}

// NewSOCResolverCaller creates a new read-only instance of SOCResolver, bound to a specific deployed contract.
func NewSOCResolverCaller(address common.Address, caller bind.ContractCaller) (*SOCResolverCaller, error) {
	contract, err := bindSOCResolver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SOCResolverCaller{contract: contract}, nil
}

// NewSOCResolverTransactor creates a new write-only instance of SOCResolver, bound to a specific deployed contract.
func NewSOCResolverTransactor(address common.Address, transactor bind.ContractTransactor) (*SOCResolverTransactor, error) {
	contract, err := bindSOCResolver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SOCResolverTransactor{contract: contract}, nil
}

// NewSOCResolverFilterer creates a new log filterer instance of SOCResolver, bound to a specific deployed contract.
func NewSOCResolverFilterer(address common.Address, filterer bind.ContractFilterer) (*SOCResolverFilterer, error) {
	contract, err := bindSOCResolver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SOCResolverFilterer{contract: contract}, nil
}

// bindSOCResolver binds a generic wrapper to an already deployed contract.
func bindSOCResolver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SOCResolverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SOCResolver *SOCResolverRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SOCResolver.Contract.SOCResolverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SOCResolver *SOCResolverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SOCResolver.Contract.SOCResolverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SOCResolver *SOCResolverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SOCResolver.Contract.SOCResolverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SOCResolver *SOCResolverCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SOCResolver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SOCResolver *SOCResolverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SOCResolver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SOCResolver *SOCResolverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SOCResolver.Contract.contract.Transact(opts, method, params...)
}

// Suffix is a free data retrieval call binding the contract method 0x98e77135.
//
// Solidity: function Suffix() constant returns(string)
func (_SOCResolver *SOCResolverCaller) Suffix(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SOCResolver.contract.Call(opts, out, "Suffix")
	return *ret0, err
}

// Suffix is a free data retrieval call binding the contract method 0x98e77135.
//
// Solidity: function Suffix() constant returns(string)
func (_SOCResolver *SOCResolverSession) Suffix() (string, error) {
	return _SOCResolver.Contract.Suffix(&_SOCResolver.CallOpts)
}

// Suffix is a free data retrieval call binding the contract method 0x98e77135.
//
// Solidity: function Suffix() constant returns(string)
func (_SOCResolver *SOCResolverCallerSession) Suffix() (string, error) {
	return _SOCResolver.Contract.Suffix(&_SOCResolver.CallOpts)
}

// GetAddrByAlias is a free data retrieval call binding the contract method 0x351af6a0.
//
// Solidity: function getAddrByAlias(_aliasName string) constant returns(address)
func (_SOCResolver *SOCResolverCaller) GetAddrByAlias(opts *bind.CallOpts, _aliasName string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SOCResolver.contract.Call(opts, out, "getAddrByAlias", _aliasName)
	return *ret0, err
}

// GetAddrByAlias is a free data retrieval call binding the contract method 0x351af6a0.
//
// Solidity: function getAddrByAlias(_aliasName string) constant returns(address)
func (_SOCResolver *SOCResolverSession) GetAddrByAlias(_aliasName string) (common.Address, error) {
	return _SOCResolver.Contract.GetAddrByAlias(&_SOCResolver.CallOpts, _aliasName)
}

// GetAddrByAlias is a free data retrieval call binding the contract method 0x351af6a0.
//
// Solidity: function getAddrByAlias(_aliasName string) constant returns(address)
func (_SOCResolver *SOCResolverCallerSession) GetAddrByAlias(_aliasName string) (common.Address, error) {
	return _SOCResolver.Contract.GetAddrByAlias(&_SOCResolver.CallOpts, _aliasName)
}

// GetAliasByAddr is a free data retrieval call binding the contract method 0x51119ae1.
//
// Solidity: function getAliasByAddr(_addr address) constant returns(string)
func (_SOCResolver *SOCResolverCaller) GetAliasByAddr(opts *bind.CallOpts, _addr common.Address) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SOCResolver.contract.Call(opts, out, "getAliasByAddr", _addr)
	return *ret0, err
}

// GetAliasByAddr is a free data retrieval call binding the contract method 0x51119ae1.
//
// Solidity: function getAliasByAddr(_addr address) constant returns(string)
func (_SOCResolver *SOCResolverSession) GetAliasByAddr(_addr common.Address) (string, error) {
	return _SOCResolver.Contract.GetAliasByAddr(&_SOCResolver.CallOpts, _addr)
}

// GetAliasByAddr is a free data retrieval call binding the contract method 0x51119ae1.
//
// Solidity: function getAliasByAddr(_addr address) constant returns(string)
func (_SOCResolver *SOCResolverCallerSession) GetAliasByAddr(_addr common.Address) (string, error) {
	return _SOCResolver.Contract.GetAliasByAddr(&_SOCResolver.CallOpts, _addr)
}

// GetInterface is a free data retrieval call binding the contract method 0xdf1827df.
//
// Solidity: function getInterface() constant returns(bytes4[])
func (_SOCResolver *SOCResolverCaller) GetInterface(opts *bind.CallOpts) ([][4]byte, error) {
	var (
		ret0 = new([][4]byte)
	)
	out := ret0
	err := _SOCResolver.contract.Call(opts, out, "getInterface")
	return *ret0, err
}

// GetInterface is a free data retrieval call binding the contract method 0xdf1827df.
//
// Solidity: function getInterface() constant returns(bytes4[])
func (_SOCResolver *SOCResolverSession) GetInterface() ([][4]byte, error) {
	return _SOCResolver.Contract.GetInterface(&_SOCResolver.CallOpts)
}

// GetInterface is a free data retrieval call binding the contract method 0xdf1827df.
//
// Solidity: function getInterface() constant returns(bytes4[])
func (_SOCResolver *SOCResolverCallerSession) GetInterface() ([][4]byte, error) {
	return _SOCResolver.Contract.GetInterface(&_SOCResolver.CallOpts)
}

// GetSuffix is a free data retrieval call binding the contract method 0x1c0e1d09.
//
// Solidity: function getSuffix() constant returns(string)
func (_SOCResolver *SOCResolverCaller) GetSuffix(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SOCResolver.contract.Call(opts, out, "getSuffix")
	return *ret0, err
}

// GetSuffix is a free data retrieval call binding the contract method 0x1c0e1d09.
//
// Solidity: function getSuffix() constant returns(string)
func (_SOCResolver *SOCResolverSession) GetSuffix() (string, error) {
	return _SOCResolver.Contract.GetSuffix(&_SOCResolver.CallOpts)
}

// GetSuffix is a free data retrieval call binding the contract method 0x1c0e1d09.
//
// Solidity: function getSuffix() constant returns(string)
func (_SOCResolver *SOCResolverCallerSession) GetSuffix() (string, error) {
	return _SOCResolver.Contract.GetSuffix(&_SOCResolver.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SOCResolver *SOCResolverCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SOCResolver.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SOCResolver *SOCResolverSession) Owner() (common.Address, error) {
	return _SOCResolver.Contract.Owner(&_SOCResolver.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SOCResolver *SOCResolverCallerSession) Owner() (common.Address, error) {
	return _SOCResolver.Contract.Owner(&_SOCResolver.CallOpts)
}

// SupportSNSInterface is a free data retrieval call binding the contract method 0xc67b08d2.
//
// Solidity: function supportSNSInterface(_interfaceId bytes4) constant returns(bool)
func (_SOCResolver *SOCResolverCaller) SupportSNSInterface(opts *bind.CallOpts, _interfaceId [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SOCResolver.contract.Call(opts, out, "supportSNSInterface", _interfaceId)
	return *ret0, err
}

// SupportSNSInterface is a free data retrieval call binding the contract method 0xc67b08d2.
//
// Solidity: function supportSNSInterface(_interfaceId bytes4) constant returns(bool)
func (_SOCResolver *SOCResolverSession) SupportSNSInterface(_interfaceId [4]byte) (bool, error) {
	return _SOCResolver.Contract.SupportSNSInterface(&_SOCResolver.CallOpts, _interfaceId)
}

// SupportSNSInterface is a free data retrieval call binding the contract method 0xc67b08d2.
//
// Solidity: function supportSNSInterface(_interfaceId bytes4) constant returns(bool)
func (_SOCResolver *SOCResolverCallerSession) SupportSNSInterface(_interfaceId [4]byte) (bool, error) {
	return _SOCResolver.Contract.SupportSNSInterface(&_SOCResolver.CallOpts, _interfaceId)
}

// SetAlias is a paid mutator transaction binding the contract method 0x3912718c.
//
// Solidity: function setAlias(_aliasName string, _aliasAddr address) returns()
func (_SOCResolver *SOCResolverTransactor) SetAlias(opts *bind.TransactOpts, _aliasName string, _aliasAddr common.Address) (*types.Transaction, error) {
	return _SOCResolver.contract.Transact(opts, "setAlias", _aliasName, _aliasAddr)
}

// SetAlias is a paid mutator transaction binding the contract method 0x3912718c.
//
// Solidity: function setAlias(_aliasName string, _aliasAddr address) returns()
func (_SOCResolver *SOCResolverSession) SetAlias(_aliasName string, _aliasAddr common.Address) (*types.Transaction, error) {
	return _SOCResolver.Contract.SetAlias(&_SOCResolver.TransactOpts, _aliasName, _aliasAddr)
}

// SetAlias is a paid mutator transaction binding the contract method 0x3912718c.
//
// Solidity: function setAlias(_aliasName string, _aliasAddr address) returns()
func (_SOCResolver *SOCResolverTransactorSession) SetAlias(_aliasName string, _aliasAddr common.Address) (*types.Transaction, error) {
	return _SOCResolver.Contract.SetAlias(&_SOCResolver.TransactOpts, _aliasName, _aliasAddr)
}

// SetVipAlias is a paid mutator transaction binding the contract method 0x45c26852.
//
// Solidity: function setVipAlias(_aliasName string, _aliasAddr address) returns()
func (_SOCResolver *SOCResolverTransactor) SetVipAlias(opts *bind.TransactOpts, _aliasName string, _aliasAddr common.Address) (*types.Transaction, error) {
	return _SOCResolver.contract.Transact(opts, "setVipAlias", _aliasName, _aliasAddr)
}

// SetVipAlias is a paid mutator transaction binding the contract method 0x45c26852.
//
// Solidity: function setVipAlias(_aliasName string, _aliasAddr address) returns()
func (_SOCResolver *SOCResolverSession) SetVipAlias(_aliasName string, _aliasAddr common.Address) (*types.Transaction, error) {
	return _SOCResolver.Contract.SetVipAlias(&_SOCResolver.TransactOpts, _aliasName, _aliasAddr)
}

// SetVipAlias is a paid mutator transaction binding the contract method 0x45c26852.
//
// Solidity: function setVipAlias(_aliasName string, _aliasAddr address) returns()
func (_SOCResolver *SOCResolverTransactorSession) SetVipAlias(_aliasName string, _aliasAddr common.Address) (*types.Transaction, error) {
	return _SOCResolver.Contract.SetVipAlias(&_SOCResolver.TransactOpts, _aliasName, _aliasAddr)
}

// WithDraw is a paid mutator transaction binding the contract method 0x0a67d2c7.
//
// Solidity: function withDraw(addr address) returns()
func (_SOCResolver *SOCResolverTransactor) WithDraw(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _SOCResolver.contract.Transact(opts, "withDraw", addr)
}

// WithDraw is a paid mutator transaction binding the contract method 0x0a67d2c7.
//
// Solidity: function withDraw(addr address) returns()
func (_SOCResolver *SOCResolverSession) WithDraw(addr common.Address) (*types.Transaction, error) {
	return _SOCResolver.Contract.WithDraw(&_SOCResolver.TransactOpts, addr)
}

// WithDraw is a paid mutator transaction binding the contract method 0x0a67d2c7.
//
// Solidity: function withDraw(addr address) returns()
func (_SOCResolver *SOCResolverTransactorSession) WithDraw(addr common.Address) (*types.Transaction, error) {
	return _SOCResolver.Contract.WithDraw(&_SOCResolver.TransactOpts, addr)
}
