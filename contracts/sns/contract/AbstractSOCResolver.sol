pragma solidity ^0.4.23;

contract AbstractSOCResolver {

    function getSuffix() public view returns(string memory);

    function setVipAlias(string memory _aliasName, address _aliasAddr) public;

    function setAlias(string memory _aliasName, address _aliasAddr) public payable;

    function getAddrByAlias(string memory _aliasName) public view returns(address);

    function getAliasByAddr(address _addr) public view returns(string memory);

    function supportSNSInterface(bytes4 _interfaceId) public pure returns(bool);

    function getInterface() public pure returns(bytes4[] memory);

    function withDraw(address addr) public;
}
