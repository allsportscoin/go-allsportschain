pragma solidity ^0.4.16;

contract owned {
    address public owner;
    constructor() public {
        owner = msg.sender;
    }
    modifier onlyOwner {
        require(msg.sender == owner);
        _;
    }
}

contract restriction {
    modifier lengthLimit(string name) {
        bytes memory strBytes = bytes(name);
        require(strBytes.length>6, "domain length too short");
        _;
    }

    modifier validChar(string name) {

        _;
    }
}

contract stringLib {

    function bytes32ToString(bytes32 x) internal pure returns (string) {
        bytes memory bytesString = new bytes(32);
        uint charCount = 0;
        for (uint j = 0; j < 32; j++) {
            byte char = byte(bytes32(uint(x) * 2 ** (8 * j)));
            if (char != 0) {
                bytesString[charCount] = char;
                charCount++;
            }
        }
        bytes memory bytesStringTrimmed = new bytes(charCount);
        for (j = 0; j < charCount; j++) {
            bytesStringTrimmed[j] = bytesString[j];
        }
        return string(bytesStringTrimmed);
    }

    function bytes32ArrayToString(bytes32[] data) internal pure returns (string) {
        bytes memory bytesString = new bytes(data.length * 32);
        uint urlLength;
        for (uint i = 0; i< data.length; i++) {
            for (uint j = 0; j < 32; j++) {
                byte char = byte(bytes32(uint(data[i]) * 2 ** (8 * j)));
                if (char != 0) {
                    bytesString[urlLength] = char;
                    urlLength += 1;
                }
            }
        }
        bytes memory bytesStringTrimmed = new bytes(urlLength);
        for (i = 0; i < urlLength; i++) {
            bytesStringTrimmed[i] = bytesString[i];
        }
        return string(bytesStringTrimmed);
    }
}

contract SNS is owned, restriction, stringLib {
    bytes4[] public RESOLVER_INTERFACE_ID = [bytes4(0x3912718c),bytes4(0x351af6a0),bytes4(0x51119ae1),bytes4(0xc67b08d2)];
    bytes4 SNS_ID = 0xef45283a;

    struct Suffix {
        string suffix;
        address addr;
        bool isAvaliable;
        string content;
        // uint256 fee;
    }

    mapping(string=>Suffix) internal Suffixes;

    string[] Names;

    mapping(address=>string) internal AddressToSuffix;

    event AddSuffix(string _suffix, address _address);

    constructor() public {

    }

    function addSuffix(string _suffixName, address _addr) public onlyOwner validChar(_suffixName) {
        require(!Suffixes[_suffixName].isAvaliable, "suffix already used");
        // require(AddressToSuffix[_addr]>0,"");
        require(bytes(AddressToSuffix[_addr]).length==0, "address already bind to other suffix resolver");

        require(supportSNSInterface(_addr), "resolver contract not support sns standard");
        Names.push(_suffixName);
        Suffixes[_suffixName] = Suffix(_suffixName, _addr, true, "");
        AddressToSuffix[_addr] = _suffixName;
        emit AddSuffix(_suffixName, _addr);
    }

    function getSuffixByName(string _suffixName) public constant returns(address _addr,
        bool _isAvaliable,
        string _content) {
        require(Suffixes[_suffixName].isAvaliable, "suffix not exist");
        Suffix memory _suffix = Suffixes[_suffixName];
        _addr = _suffix.addr;
        _isAvaliable = _suffix.isAvaliable;
        _content = _suffix.content;
    }

    function getSuffixByAddr(address _addr) public constant returns(string _suffixName,
        bool _isAvaliable,
        string _content) {
        require(bytes(AddressToSuffix[_addr]).length>0, "address not bind to any suffix");
        _suffixName = AddressToSuffix[_addr];
        Suffix memory _suffix = Suffixes[_suffixName];
        _isAvaliable = _suffix.isAvaliable;
        _content = _suffix.content;
    }

    function setAliasOfSuffix(string _suffixName, string _aliasName, address _aliasAddr) public payable {
        require(Suffixes[_suffixName].isAvaliable, "suffix not exist");
        Suffix memory _suffix = Suffixes[_suffixName];
        Resolver _resol = Resolver(_suffix.addr);
        uint256 fee = msg.value;
        _resol.setAlias.value(fee)(_aliasName, _aliasAddr);
    }

    function getAliasByAddrAndSuffix(string _suffixName, address _aliasAddr) public constant returns(string) {
        require(Suffixes[_suffixName].isAvaliable, "suffix not exist");
        Suffix memory _suffix = Suffixes[_suffixName];
        Resolver _resol = Resolver(_suffix.addr);
        return _resol.getAliasByAddr(_aliasAddr);
    }

    function getAddrByAliasAndSuffix(string _suffixName, string _aliasName) public constant returns(address) {
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
            _resol.supportSNSInterface(RESOLVER_INTERFACE_ID[3]);
    }

    function withDraw(address addr) public onlyOwner {
        require(addr != 0x0);
        addr.transfer(address(this).balance);
    }
}

interface Resolver {

    function setAlias(string _aliasName, address _aliasAddr) public payable;

    function getAddrByAlias(string _aliasName) public view returns(address);

    function getAliasByAddr(address _addr) public view returns(string);

    function supportSNSInterface(bytes4 _interfaceId) public view returns(bool);

}


