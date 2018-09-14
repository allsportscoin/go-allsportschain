pragma solidity ^0.4.16;

contract Owner{
     address internal owner;

     modifier onlyOwner() {
        require(owner == msg.sender );
        _;
    }
}

contract MatchData is Owner{
    //data
    struct Item {
        string name;
        uint256 price;
        string desc;
        bool state;
        uint8 matchstate;
        uint8 ascore;
        uint8 bscore;
    }
    
    Item[] public items;
    
    //function
    constructor () payable public{
        owner = msg.sender;
    }
    
    function GetMatchInfo(uint8 id) public view returns (string name, uint256 pr, string desc, bool state){
        require(id < items.length);
        name = items[id].name;
        pr = items[id].price;
        desc = items[id].desc;
        state = items[id].state;
    }
    
    function GetMatchMaxId() public view returns (uint256) {
        return items.length;
    }
    
    function SetMatch(
        uint8 id, 
        string name, 
        uint256 pr, 
        string desc, 
        bool state,
        uint8 matchstate,
        uint8 ascore,
        uint8 bscore) public onlyOwner returns (bool)
    {
        require(id < items.length);

        items[id].name = name;
        items[id].price = pr;
        items[id].desc = desc;
        items[id].state = state;
        items[id].matchstate = matchstate;
        items[id].ascore = ascore;
        items[id].bscore = bscore;
        return true;
    }
    
    function SetMatchScore(
        uint8 id, 
        uint8 matchstate,
        uint8 ascore,
        uint8 bscore) public onlyOwner returns (bool)
    {
        require(id < items.length);
        items[id].matchstate = matchstate;
        items[id].ascore = ascore;
        items[id].bscore = bscore;
        return true;
    }
    
    function AddMatch(
        string name, 
        uint256 pr, 
        string desc, 
        bool state,
        uint8 matchstate) public onlyOwner returns (uint256)
    {
        items.push(Item(name,pr,desc,state,matchstate,0,0));
        return items.length-1;
    }
    
    function SetMatchState(uint8 id, bool state) public onlyOwner returns (bool)
    {
        require(id < items.length);
        items[id].state = state;
        return state;
    }
    
    function SetMatchMstate(uint8 id, uint8 matchstate) public onlyOwner returns (uint8)
    {
        require(id < items.length);
        items[id].matchstate = matchstate;
        return matchstate;
    }
    
    function SetMatchPrice(uint8 id, uint256 pr) public onlyOwner {
        require(id < items.length);
        items[id].price = pr;
    }
    
    function GetMatchPrice(uint8 id) public view returns (uint256){
        require(id < items.length);
        return items[id].price;
    }
    
    function GetMatchData(uint8 id) public payable returns  (uint8 matchstate, uint8 ascore, uint8 bscore){
        require(id < items.length);
        require(msg.value >= items[id].price);
        require(true >= items[id].state);
        
        //name = items[id].name;
        matchstate = items[id].matchstate;
        ascore = items[id].ascore;
        bscore = items[id].bscore;
        
        if (msg.value > items[id].price) {
            address(msg.sender).transfer(msg.value - items[id].price);
        }
    }
    
    function GetMatchByOwner(uint8 id) public view onlyOwner 
    returns (
        string name, 
        uint256 pr, 
        string desc, 
        bool state,
        uint8 matchstate,
        uint8 ascore,
        uint8 bscore
        ){
        require(id < items.length);
        name = items[id].name;
        pr = items[id].price;
        desc = items[id].desc;
        state = items[id].state;
        matchstate = items[id].matchstate;
        ascore = items[id].ascore;
        bscore = items[id].bscore;
    }
    
    function GetBalance() public view returns (uint256) {
        return address(this).balance;
    }
    

    function Draw(address drawer) public onlyOwner{
        drawer.transfer(address(this).balance);     
    }
    
    function () public payable {
    }
}

contract ScoreGuess is Owner{

    //data
    struct Player {
        address addr;
        uint256 betvalue;
    }
    
    mapping (uint256 => Player[]) nAWinB;
    mapping (uint256 => mapping (address => uint256)) nAWinBIndex;
    mapping (uint256 => Player[]) nBWinA;
    mapping (uint256 => mapping (address => uint256)) nBWinAIndex;
    mapping (uint256 => Player[]) nAEqualB;
    mapping (uint256 => mapping (address => uint256)) nAEqualBIndex;
        
    struct Guessmatch {
        string name;
        address dataAddr;
        uint8 matchid;
        uint256 totaltoken;
        uint256 nAWinBtoken;
        uint256 nBWinAtoken;
        uint256 nAEqualBtoken;
        uint8 state;
    }
    
    Guessmatch[] public matchs;
    
    //modifier
    
    modifier isHuman() {
        address _addr = msg.sender;
        uint256 _codeLength;
        
        assembly {_codeLength := extcodesize(_addr)}
        require(_codeLength == 0, "sorry humans only");
        _;
    }
    
    //function
    constructor () payable public{
        
    }
    
    function CreateMatch (string name, address Addr, uint8 matchid) public returns (uint256) {
        matchs.push(Guessmatch(name,Addr,matchid,0,0,0,0,0));
        return matchs.length - 1;
    }
    
    function BetMatch (uint256 matchid, uint8 btype) public payable returns (bool) {
        require(msg.value > 0);
        require(btype >= 0 && btype < 3);
        require(matchid < matchs.length);
        address addr = msg.sender;
        uint256 betvalue = msg.value;
        
        uint256 Index = 0;
        if (btype == 1) {//A Win B
            Index = nAWinBIndex[matchid][addr];
            if(Index == 0) {
                nAWinB[matchid].push(Player(addr, betvalue));
                nAWinBIndex[matchid][addr] = [matchid].length;
            } else {
                nAWinB[matchid][Index].betvalue += betvalue;
            }
            matchs[matchid].nAWinBtoken += betvalue;
        } else if (btype == 2){// B Win A
            Index = nBWinAIndex[matchid][addr];
            if(Index == 0) {
                nBWinA[matchid].push(Player(addr, betvalue));
                nBWinAIndex[matchid][addr] = nBWinA[matchid].length;
            } else {
                nBWinA[matchid][Index].betvalue += betvalue;
            }
            matchs[matchid].nBWinAtoken += betvalue;
        } else {// A equal B
            Index = nAEqualBIndex[matchid][addr];
            if(Index == 0) {
                nAEqualB[matchid].push(Player(addr, betvalue));
                nAEqualBIndex[matchid][addr] = nAEqualB[matchid].length;
            } else {
                nAEqualB[matchid][Index].betvalue += betvalue;
            }
            matchs[matchid].nAEqualBtoken += betvalue;
        }
        matchs[matchid].totaltoken += betvalue;
        return true;
    }
    
    function transferToken (
        Player[] user, 
        uint256 pooltoken, 
        uint256 alltotal) internal returns (bool) {
            
        if (user.length <= 0 
            || pooltoken < user.length
            || alltotal < user.length) {
            return false;
        }
        
        uint8 i = 0;
        for (i=0; i<user.length; i++) {
            address(user[i].addr).transfer(user[i].betvalue/pooltoken*alltotal);
        }
        return true;
    }
    
    function draw(uint256 matchid) public payable returns (bool) {
        MatchData MD = MatchData(matchs[matchid].dataAddr);
        uint8 matchstate;
        uint8 ascore;
        uint8 bscore;
        (matchstate, ascore, bscore) = MD.GetMatchData.value(1)(matchs[matchid].matchid);
        if (matchstate != 1 || matchs[matchid].state != 0) {
            return false;
        }
        
        if (ascore > bscore) {// A Win B
            transferToken(nAWinB[matchid], matchs[matchid].nAWinBtoken, matchs[matchid].totaltoken);
        } else if (ascore < bscore) {// B Win A
            transferToken(nBWinA[matchid], matchs[matchid].nBWinAtoken, matchs[matchid].totaltoken);
        } else {// A Equal B
            transferToken(nAEqualB[matchid], matchs[matchid].nAEqualBtoken, matchs[matchid].totaltoken);
        }
        
        matchs[matchid].state = 1;
        
        return true;
    }
    
    function Pay () public payable {
    }
    
    function GetBalance() public view returns(uint256){
        return address(this).balance;
    }
    
    
    function () public payable {
    }
    
}
