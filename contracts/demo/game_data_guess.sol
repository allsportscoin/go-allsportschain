pragma solidity ^0.4.17;


contract Owner{
    address public owner;

    modifier onlyOwner() {
        if(msg.sender != owner) revert();
        else _;
    }

}

contract GameData is Owner
{
    struct Game {
        uint256 id;
        string  name;
        uint256 startTime;
        uint256 endTime;
        uint8 state;//0:未开始，1:进行中 2: 结束
        string teamA;
        string teamB;
        uint8 teamAScore;
        uint8 teamBScore;
    }

    Game[] public gameList ;
    address admin;

    string public name;
    string public symbol;
    uint256 public fee;

    modifier onlyAdmin()
    {
        if(msg.sender != admin) revert();
        else _;
    }
    modifier needFee()
    {
        if(msg.value < fee ) revert();
        else _;
    }

    function() public { revert(); } //should not send any ether directly

    //"mdd GameData","GameData","0xca35b7d915458ef540ade6068dfe2f44e8fa733c"
    constructor(
        string glName,
        string glSymbol,
        address adminAddr) public {
        //assign some tickets to event admin
        admin = adminAddr;
        symbol = glSymbol;
        name = glName;
        fee = 1;
    }

    function getBalance() public view returns(uint256){
        address myAddress = this;
        return myAddress.balance;
    }

    function setFee(uint256 value) public onlyOwner returns(bool){
        fee = value;
        return true;
    }

    function name() public view returns(string)
    {
        return name;
    }

    function symbol() public view returns(string)
    {
        return symbol;
    }

    function getGameListSize() public view returns(uint256 size)
    {
        return gameList.length;
    }

    //"himdd",100,200,1,"A","B",1,2
    function addGame(string gName, uint256 gStartTime,uint256 gEndTime,uint8 gState, string gTeamA, string gTeamB, uint8  gTeamAScore,uint8 gTeamBScore) onlyOwner public
    {
        uint256 gId = gameList.length;
        Game memory game = Game(gId, gName, gStartTime, gEndTime, gState, gTeamA, gTeamB, gTeamAScore, gTeamBScore);
        gameList.push(game);
    }

    function setGameState(uint256 needIndex,uint8 gState) public onlyOwner {
        gameList[needIndex].state = gState;
    }

    function getGameState(uint256 needIndex) public view returns(uint8 gState){
        return gameList[needIndex].state;
    }

    function setGameScore(uint256 needIndex,uint8 gTeamAScore,uint8 gTeamBScore) public onlyOwner {
        gameList[needIndex].teamAScore = gTeamAScore;
        gameList[needIndex].teamBScore = gTeamBScore;
    }

    function getGameScore(uint256 needIndex) public payable needFee returns(uint8  gTeamAScore,uint8 gTeamBScore){
        Game memory game = gameList[needIndex];
        return (game.teamAScore,game.teamBScore);
    }

    function getGame(uint256 needIndex) public payable needFee returns(string gName, uint256 gStartTime,uint256 gEndTime,uint8 gState, string gTeamA, string gTeamB, uint8  gTeamAScore,uint8 gTeamBScore){
        Game memory game = gameList[needIndex];
        return (game.name,game.startTime,game.endTime,game.state,game.teamA,game.teamB,game.teamAScore,game.teamBScore);
    }

    function getGameList(uint256 uState) public view returns(uint256[] list,uint256 size){
        uint256[] memory arr = new uint256[](gameList.length);
        uint256 tIndex = 0;
        for (uint256 i = 0; i< gameList.length; i++){
            Game memory game = gameList[i];
            if(game.state  == uState){
                arr[tIndex++] = i;
            }
        }
        return (arr, tIndex);
    }

    function getGameByOwner(uint256 needIndex) public view onlyOwner returns(string gName, uint256 gStartTime,uint256 gEndTime,uint8 gState, string gTeamA, string gTeamB, uint8  gTeamAScore,uint8 gTeamBScore){
        Game memory game = gameList[needIndex];
        return (game.name,game.startTime,game.endTime,game.state,game.teamA,game.teamB,game.teamAScore,game.teamBScore);
    }

    function getGameListByOwner(uint256 uState) public view onlyOwner returns(uint256[] list,uint256 size){
        uint256[] memory arr = new uint256[]( gameList.length);
        uint256 tIndex = 0;
        for (uint256 i = 0; i< gameList.length; i++){
            Game memory game = gameList[i];
            if(game.state  == uState){
                arr[tIndex++] = i;
            }
        }
        return (arr, tIndex);
    }


    function endContract() onlyOwner public
    {
        selfdestruct(owner);
    }

    function getContractAddress() public view returns(address)
    {
        return this;
    }

}

contract GuessApp is Owner {
    uint256 public fee;

    struct Guess{
        uint256 id;
        uint256 gameId;
        uint256 totalInPool;
        uint256 winTotal;
        uint256 loseTotal;
        uint256 deuceTotal;
    }

    Guess[] guessList;

    mapping(uint256=> mapping (address => uint256)) winMap;
    mapping(uint256=> mapping (address => uint256)) loseMap;
    mapping(uint256=> mapping (address => uint256)) deuceMap;
    mapping(uint256=> address[]) winUsersMap;
    mapping(uint256=> address[]) loseUsersMap;
    mapping(uint256=> address[]) deuceUsersMap;

    constructor() public payable {
        fee = 1;
    }

    function setFee(uint256 value) public onlyOwner returns(bool){
        fee = value;
        return true;
    }
    function withDraw(address to, uint256 value) public onlyOwner {
        to.transfer(value);
    }

    function getBalance() public view returns(uint256){
        address myAddress = this;
        return myAddress.balance;
    }

    function ()  public payable {
    }

    function addGuess(uint256 gameId) public returns(bool){
        guessList.push(Guess(guessList.length,gameId,0,0,0,0));
        return true;
    }

    function chipIn(uint256 guessId, uint8 chipType) public payable {
        if (chipType == 1){
            winUsersMap[guessId].push(msg.sender);
            winMap[guessId][msg.sender] += msg.value;
            guessList[guessId].winTotal += msg.value;
            guessList[guessId].totalInPool += msg.value;
        }else if (chipType == 2){
            loseUsersMap[guessId].push(msg.sender);
            loseMap[guessId][msg.sender] += msg.value;
            guessList[guessId].loseTotal += msg.value;
            guessList[guessId].totalInPool += msg.value;
        }else{
            deuceUsersMap[guessId].push(msg.sender);
            deuceMap[guessId][msg.sender] += msg.value;
            guessList[guessId].deuceTotal += msg.value;
            guessList[guessId].totalInPool += msg.value;
        }
    }

    function lottery(address contractAddr, uint256 guessId) public payable returns(bool){
        GameData gameData = GameData(contractAddr);
        uint256 gameId = guessList[guessId].gameId;
        uint8 gState = gameData.getGameState(gameId);
        if(gState != 2){
            return false;
        }
        uint8 teamAScore;
        uint8 teamBScore;
        ( teamAScore,teamBScore) = gameData.getGameScore.value(fee)(gameId);
        address addr ;
        uint256 coin;
        uint256 j;
        if( teamAScore > teamBScore ){
            for(  j=0; j< winUsersMap[guessId].length; j++){
                addr = winUsersMap[guessId][j];
                coin = guessList[guessId].totalInPool * winMap[guessId][addr] /guessList[guessId].winTotal;
                addr.transfer(coin);
            }

        }else if(teamAScore < teamBScore ){
            for( j=0; j< loseUsersMap[guessId].length; j++){
                addr = loseUsersMap[guessId][j];
                coin = guessList[guessId].totalInPool * loseMap[guessId][addr] /guessList[guessId].loseTotal;
                addr.transfer(coin);
            }
        }else{
            for( j=0; j< deuceUsersMap[guessId].length; j++){
                addr = deuceUsersMap[guessId][j];
                coin = guessList[guessId].totalInPool * deuceMap[guessId][addr] /guessList[guessId].deuceTotal;
                addr.transfer(coin);
            }

        }

    }


    function callGetGameState(address contractAddr,uint256 gameId) public payable returns(uint8 gState){
        GameData gameData = GameData(contractAddr);
        return gameData.getGameState(gameId);
    }

    function callGetGameScore(address contractAddr,uint256 gameId) public payable returns(uint8  gTeamAScore,uint8 gTeamBScore){
        GameData gameData = GameData(contractAddr);
        return gameData.getGameScore.value(fee)(gameId);
    }

    function callGetGame(address contractAddr,uint256 gameId) public payable returns(string gName, uint256 gStartTime,uint256 gEndTime,uint8 gState,string gTeamA, string gTeamB, uint8  gTeamAScore,uint8 gTeamBScore){
        GameData gameData = GameData(contractAddr);
        return gameData.getGame.value(fee)(gameId);
    }
    function callGetGameList(address contractAddr,uint256 uState) public view returns(uint256[] list,uint256 size){
        GameData gameData = GameData(contractAddr);
        return gameData.getGameList(uState);
    }
}