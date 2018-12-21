pragma solidity ^0.4.23;

contract AbstractSNS {

    function equals(string memory str1, string memory str2) internal pure returns(bool);

    function addSuffix(string memory _suffixName, address _addr) public;

    function getSuffixByName(string memory _suffixName) public view returns(address _addr,
        bool _isAvaliable,
        string memory _content);

    function getSuffixByAddr(address _addr) public view returns(string memory _suffixName,
        bool _isAvaliable,
        string memory _content);

    function setAliasOfSuffix(string memory _suffixName, string memory _aliasName, address _aliasAddr) public payable;

    function getAliasByAddrAndSuffix(string memory _suffixName, address _aliasAddr) public view returns(string memory);

    function getAddrByAliasAndSuffix(string memory _suffixName, string memory _aliasName) public view returns(address);

    function supportSNSInterface(address _addr) public view returns(bool);

    function withDraw(address addr) public;
}
