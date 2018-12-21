pragma solidity ^0.4.23;

import "./AbstractSNS.sol";

contract SNS is AbstractSNS {
    bytes4[] public RESOLVER_INTERFACE_ID = [bytes4(0x3912718c),bytes4(0x351af6a0),bytes4(0x51119ae1),bytes4(0xc67b08d2),bytes4(0x1c0e1d09)];
    // bytes4 SNS_ID = 0xef45283a;
    address public owner;

    string[] Names;
    mapping(string=>Suffix) internal Suffixes;
    mapping(address=>string) internal AddressToSuffix;

    struct Suffix {
        string suffix;
        address addr;
        bool isAvaliable;
        string content;
        // uint256 fee;
    }

    // event AddSuffix(string _suffix, address _address);

    modifier onlyOwner {
        require(msg.sender == owner);
        _;
    }

    //check suffix validation:length/suffix format(like .soc)
    modifier validSuffix(string memory _suffix) {
        bytes memory _suffixByte = bytes(_suffix);
        uint suffixLength = _suffixByte.length;

        require((suffixLength>2 && suffixLength<6), "length of suffix include dot. keeps 3-6");

        require(_suffixByte[0]>=bytes32(".") ,
                "Initial character must be dot(.)."
        );

        for(uint i=1; i<suffixLength; i++) {
            require((_suffixByte[i]>=bytes32("a") && _suffixByte[i]<=bytes32("z")),
                "invalid character error."
            );
        }
        _;
    }

    function equals(string memory str1, string memory str2) internal pure returns(bool) {
        bytes memory byte1 = bytes(str1);
        bytes memory byte2 = bytes(str2);
        require(byte1.length == byte2.length, "strings not equals");
        uint strLength = byte1.length;
        for(uint i=0; i<strLength; i++) {
            require(byte1[i] == byte2[i], "strings not equals");
        }
        return true;
    }

    constructor() public {
        owner = msg.sender;
    }

    function addSuffix(string memory _suffixName, address _addr) public onlyOwner validSuffix(_suffixName) {
        require(!Suffixes[_suffixName].isAvaliable, "suffix already used");
        // require(AddressToSuffix[_addr]>0,"");
        require(bytes(AddressToSuffix[_addr]).length==0, "address already bind to other suffix resolver");

        Resolver _resol = Resolver(_addr);
        string memory suffixName = _resol.getSuffix();
        require(equals(_suffixName, suffixName), "suffix name invalid");

        require(supportSNSInterface(_addr), "resolver contract not support sns standard");
        Names.push(_suffixName);
        Suffixes[_suffixName] = Suffix(_suffixName, _addr, true, "");
        AddressToSuffix[_addr] = _suffixName;
        // emit AddSuffix(_suffixName, _addr);
    }

    function getSuffixByName(string memory _suffixName) public view returns(address _addr,
        bool _isAvaliable,
        string memory _content) {
        require(Suffixes[_suffixName].isAvaliable, "suffix not exist");
        Suffix memory _suffix = Suffixes[_suffixName];
        _addr = _suffix.addr;
        _isAvaliable = _suffix.isAvaliable;
        _content = _suffix.content;
    }

    function getSuffixByAddr(address _addr) public view returns(string memory _suffixName,
        bool _isAvaliable,
        string memory _content) {
        require(bytes(AddressToSuffix[_addr]).length>0, "address not bind to any suffix");
        _suffixName = AddressToSuffix[_addr];
        Suffix memory _suffix = Suffixes[_suffixName];
        _isAvaliable = _suffix.isAvaliable;
        _content = _suffix.content;
    }

    function setAliasOfSuffix(string memory _suffixName, string memory _aliasName, address _aliasAddr) public payable {
        require(Suffixes[_suffixName].isAvaliable, "suffix not exist");
        Suffix memory _suffix = Suffixes[_suffixName];
        Resolver _resol = Resolver(_suffix.addr);
        uint256 fee = msg.value;
        _resol.setAlias.value(fee)(_aliasName, _aliasAddr);
    }

    function getAliasByAddrAndSuffix(string memory _suffixName, address _aliasAddr) public view returns(string memory) {
        require(Suffixes[_suffixName].isAvaliable, "suffix not exist");
        Suffix memory _suffix = Suffixes[_suffixName];
        Resolver _resol = Resolver(_suffix.addr);
        return _resol.getAliasByAddr(_aliasAddr);
    }

    function getAddrByAliasAndSuffix(string memory _suffixName, string memory _aliasName) public view returns(address) {
        require(Suffixes[_suffixName].isAvaliable, "suffix not exist");
        Suffix memory _suffix = Suffixes[_suffixName];
        Resolver _resol = Resolver(_suffix.addr);
        return _resol.getAddrByAlias(_aliasName);
    }

    function supportSNSInterface(address _addr) public view returns(bool) {
        Resolver _resol = Resolver(_addr);
        return _resol.supportSNSInterface(RESOLVER_INTERFACE_ID[0]) &&
            _resol.supportSNSInterface(RESOLVER_INTERFACE_ID[1]) &&
            _resol.supportSNSInterface(RESOLVER_INTERFACE_ID[2]) &&
            _resol.supportSNSInterface(RESOLVER_INTERFACE_ID[3]) &&
            _resol.supportSNSInterface(RESOLVER_INTERFACE_ID[4]);
    }

    function withDraw(address addr) public onlyOwner {
        require(addr != address(0x0));
        addr.transfer(address(this).balance);
    }
}

interface Resolver {

    function getSuffix() external view returns(string memory);

    function setAlias(string _aliasName, address _aliasAddr) external payable;

    function getAddrByAlias(string _aliasName) external view returns(address);

    function getAliasByAddr(address _addr) external view returns(string memory);

    function supportSNSInterface(bytes4 _interfaceId) external view returns(bool);

}


