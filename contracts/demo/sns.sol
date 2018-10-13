pragma solidity ^0.4.0;

 contract owned {
     address public owner;

     function owned() public {
         owner = msg.sender;
     }

     modifier onlyOwner {
         require(msg.sender == owner);
         _;
     }

 }

 contract stringLib {

     function bytes32ToString(bytes32 x) internal constant returns (string) {
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

     function bytes32ArrayToString(bytes32[] data) internal constant returns (string) {
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

 contract SNS is owned, stringLib{

     struct Record {
         address addr;
         string content;
         bytes32 bname;
         string name;
         string suffix;

     }

     mapping(bytes32 => Record) public Records;
     mapping(bytes32 => bytes32) internal record_index;
     mapping(bytes32 => uint) internal suffix_level; // 1: onlyOwner, >1 userable
     string[] public Suffixs;
     uint256 public fee = 1;

     event Transfer(address _from, address _to, string _name);
     event ChangeName(address _owner, string _old, string _new);
     event NewSuffix(string _suffix, uint owner);

     function SNS(uint256 _fee) payable {
         if(_fee > 0){
             fee = _fee;
         }
     }

     function changeFee(uint256 _fee) public onlyOwner {
         fee = _fee;
     }

     //for owner add some suffix for users, just like, "soc","user" etc
     function addSuffix(bytes32 _suffix, uint _ownerOnly) public onlyOwner {
         suffix_level[_suffix] = _ownerOnly;
         string memory suffix = bytes32ToString(_suffix);
         Suffixs.push(suffix);
         NewSuffix(suffix, _ownerOnly);
     }

     function getRecord(bytes32 _name) public returns(address _addr, string _suffix) {
         require(record_index[_name].length > 0, 'this name not used');
         Record record = Records[_name];
         _addr = record.addr;
         _suffix = record.suffix;
     }

     //user can bind a name and suffix with his address, when the suffix exists and the name is unused
     //Anyone who want bind name should pay some token
     //only name owner can change his domain name
     function setName(bytes32 _name, bytes32 _suffix) payable public {
         require(msg.value >= fee, 'must pay some value from bind name');
         require(suffix_level[_suffix] > 0, 'no such suffix');

         if(record_index[_name].length > 0){
             require(record.addr == msg.sender, 'no Permission Change the record');
         }
         if(suffix_level[_suffix] == 1){
             require(msg.sender == owner, 'this suffix only owner used');
         }
         Record record = Records[_name];

         string oldName = record.name;
         string memory name = bytes32ToString(_name);
         string memory suffix = bytes32ToString(_suffix);
         record.name = name;
         record.suffix = suffix;
         record.content = '';
         record.addr = msg.sender;

         record_index[_name] = _suffix;
         if(record.bname.length > 0){
             delete record_index[record.bname];
         }
         record.bname = _name;
         ChangeName(msg.sender, oldName, name);

     }

     //user can transfer his name to other
     function transferName(bytes32 _name, address _addr) public {
         require(record_index[_name].length > 0);
         Record record = Records[_name];
         require(record.addr == msg.sender, 'no Permission Change the record');

         record.addr = _addr;
         Transfer(msg.sender, _addr, bytes32ToString(_name));
     }

     //with draw
     function withDraw() onlyOwner {
         owner.transfer(address(this).balance);
     }
 }
