pragma solidity ^0.4.23;

import "./AbstractSOCResolver.sol";

contract SOCResolver is AbstractSOCResolver {
    address public owner;
    //format like ".soc"
    string public Suffix;

    string[] AliasNames;
    mapping(string=>Alias) Aliases;
    mapping(address=>string) AddressToAlias;

    // event SetAlias(string aliasName, address aliasAddr);

    struct Alias {
        string name;
        address addr;
        bool isUsed;
    }

    modifier onlyOwner {
        require(msg.sender == owner);
        _;
    }

    modifier validAlias(string memory _alias, string memory _suffix) {
        bytes memory _aliasByte = bytes(_alias);
        bytes memory _suffixByte = bytes(_suffix);
        uint aliasLength = _aliasByte.length;
        uint suffixLength = _suffixByte.length;

        require((aliasLength>6+suffixLength && aliasLength<13+suffixLength), "domain length except suffix keeps 7-12");

        for (uint i=0; i<suffixLength; i++) {
            require(_aliasByte[aliasLength-i-1] == _suffixByte[suffixLength-i-1], "invalid suffix error");
        }

        require((_aliasByte[0]>=bytes32("A") && _aliasByte[0]<=bytes32("Z"))
                || (_aliasByte[0]>=bytes32("a") && _aliasByte[0]<=bytes32("z")),
                "Initial character must be english alphabet."
        );

        for(uint j=1; j<aliasLength - suffixLength; j++) {

            require((_aliasByte[j]>=bytes32("0") && _aliasByte[j]<=bytes32("9"))
                || (_aliasByte[j]>=bytes32("A") && _aliasByte[j]<=bytes32("Z"))
                || (_aliasByte[j]>=bytes32("a") && _aliasByte[j]<=bytes32("z")),
                "invalid character error."
            );
        }
        _;
    }



    constructor(string memory _suffix) public {
        owner = msg.sender;
        Suffix = _suffix;
    }

    function setVipAlias(string memory _aliasName, address _aliasAddr) public onlyOwner {
        // require(validAlias(_aliasName, Suffix), "invalid alias format.");
        require(!Aliases[_aliasName].isUsed, "alias already exists");
        require(_aliasAddr != address(0x0), "empty address not allowed");
        require(bytes(AddressToAlias[_aliasAddr]).length==0, "address already bind to another alias");

        AliasNames.push(_aliasName);
        AddressToAlias[_aliasAddr] = _aliasName;
        Aliases[_aliasName] = Alias(_aliasName, _aliasAddr, true);

        // emit SetAlias(_aliasName, _aliasAddr);
    }

    function setAlias(string memory _aliasName, address _aliasAddr) public payable validAlias(_aliasName,Suffix) {
        require(!Aliases[_aliasName].isUsed, "alias already exists");
        if(_aliasAddr == address(0x0)) {
            _aliasAddr = msg.sender;
        }
        require(bytes(AddressToAlias[_aliasAddr]).length==0, "address already bind to another alias");

        AliasNames.push(_aliasName);
        AddressToAlias[_aliasAddr] = _aliasName;
        Aliases[_aliasName] = Alias(_aliasName, _aliasAddr, true);

        // emit SetAlias(_aliasName, _aliasAddr);
    }

    function getAddrByAlias(string memory _aliasName) public view returns(address) {
        require(Aliases[_aliasName].isUsed, "alias name not exists");
        Alias memory _alias = Aliases[_aliasName];
        return _alias.addr;
    }

    function getAliasByAddr(address _addr) public view returns(string memory) {
        require(bytes(AddressToAlias[_addr]).length>0, "alias of _addr not exists");
        Alias memory _alias = Aliases[AddressToAlias[_addr]];
        return _alias.name;
    }

    function getSuffix() public view returns(string memory) {
        return Suffix;
    }

    function supportSNSInterface(bytes4 _interfaceId) public pure returns(bool) {
        return this.getAddrByAlias.selector==_interfaceId ||
            this.getAliasByAddr.selector==_interfaceId ||
            this.setAlias.selector==_interfaceId ||
            this.supportSNSInterface.selector==_interfaceId ||
            this.getSuffix.selector==_interfaceId;
    }

    function getInterface() public pure returns(bytes4[] memory) {
        bytes4[] memory res = new bytes4[](5);
        res[0] = bytes4(keccak256("setAlias(string,address)"));
        res[1] = bytes4(keccak256("getAddrByAlias(string)"));
        res[2] = bytes4(keccak256("getAliasByAddr(address)"));
        res[3] = bytes4(keccak256("supportSNSInterface(bytes4)"));
        res[4] = bytes4(keccak256("getSuffix()"));

        return res;
    }

    function withDraw(address addr) public onlyOwner {
        require(addr != address(0x0));
        addr.transfer(address(this).balance);
    }
}
