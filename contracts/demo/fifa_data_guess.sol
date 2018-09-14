pragma solidity ^0.4.22;


contract FifaData {

    struct Match {
        uint id;
        bytes32 official_name;
        bytes32 home_team;
        bytes32 guest_team;
        uint start_time;
        uint end_time;
        //enum State {prepare, play, end}
        uint state;
        uint home_team_goals;
        uint guest_team_goals;
    }
    
    //mapping(uint => Match) matchsIdMap;
    //uint max_id;

    Match[] matchsIdMap;

    uint256 public fee;
    
    address public owner;
    
    modifier onlyOwner() {
        require(
            msg.sender == owner,
            "Only owner can call this."
        );
        _;
    }
    
    modifier feeValue() {
        require(msg.value >= fee);
        _;
    }
    
    //init contract
    constructor() public payable {
        owner = msg.sender;
        fee = 1;
    }
    
    //set fee each call
    function setFee(uint256 value) public onlyOwner  
    returns(bool)
    {
        fee = value;
        return true;
    }
    
    //get max match id
    function getMaxMatchId() public returns (uint)
    {
        return matchsIdMap.length;
    }
    
    //get time range matchs 
    function getNewMatchIdsAfter24Hours(uint time_start) 
    public returns(uint[])
    {
        uint[] memory id_list = new uint[](matchsIdMap.length);
        uint index = 0;
        for(uint i = 0; i < matchsIdMap.length; i++) {
            uint id = matchsIdMap[i].id;
            uint start_time = matchsIdMap[i].start_time;
            if(start_time - time_start > 24*60*60) {
                id_list[index++] = matchsIdMap[i].id;
            }
        }
        return id_list;
    }
    
    //del start and end time, for stack too deep error
    function getMatchById(uint id) 
    public payable feeValue
    returns (uint, bytes32, bytes32, bytes32, uint , uint, uint) 
    {
        return (matchsIdMap[id].id, matchsIdMap[id].official_name, 
        matchsIdMap[id].home_team, matchsIdMap[id].guest_team,
        //matchsIdMap[id].start_time, matchsIdMap[id].end_time,
        matchsIdMap[id].state, matchsIdMap[id].home_team_goals, matchsIdMap[id].guest_team_goals);
    }
    
    //owner check data
    function ownerGetMatchById(uint256 id) 
    public onlyOwner 
    returns (uint, bytes32, bytes32, bytes32, uint , uint, uint) 
    {
        return (matchsIdMap[id].id, matchsIdMap[id].official_name, 
        matchsIdMap[id].home_team, matchsIdMap[id].guest_team,
        //matchsIdMap[id].start_time, matchsIdMap[id].end_time,
        matchsIdMap[id].state, matchsIdMap[id].home_team_goals, matchsIdMap[id].guest_team_goals);
    }
    
    function getMatchStateAndGoals(uint256 id)
    public 
    returns (uint256, uint256, uint256)
    {
        return (matchsIdMap[id].state, matchsIdMap[id].home_team_goals, matchsIdMap[id].guest_team_goals);
    }
    
    //add new match
    function addMatch(bytes32 official_name, bytes32 team_a, bytes32 team_b, uint start_time, uint end_time) 
    public onlyOwner
    {
        uint id = matchsIdMap.length;
        Match memory new_match = Match(id, official_name, team_a, team_b, start_time, end_time, 0,0,0);
        matchsIdMap.push(new_match);
    }
    
    //set match result
    function setMatchResult(uint id, uint home_team_goals, uint guest_team_goals)
    public onlyOwner 
    {
        matchsIdMap[id].home_team_goals = home_team_goals;
        matchsIdMap[id].guest_team_goals = guest_team_goals;
    }
    
    //set mach state
    function setMatchState(uint id, uint state)
    public onlyOwner 
    {
        matchsIdMap[id].state = state;
    }
    
}


contract FootballGuess  {
    address public owner;

    mapping(uint256 => address[]) matchAddressMap;
    mapping(uint256 => uint256[3]) matchOdds;
    mapping(uint256 => uint256[3]) matchSums;
    mapping(address => uint256[][3]) addressMatchsMap;

    constructor() public {
        owner = msg.sender;    
    }

    //get match list
    function getFifaDataMatchMaxId(address contractAddr) public returns(uint256){
        FifaData data = FifaData(contractAddr);
        uint256 max_id = data.getMaxMatchId();
        return max_id;
    }
    
    //return odds and other info
    function getMatchDetailInfo(address contractAddr, uint256 matchId) public {
        FifaData data = FifaData(contractAddr);
        data.getMatchById(matchId);
        //matchOdds[matchId]
    }
    
    //add match  guess record
    function addMatchGuessRecord(uint256 matchId, uint256 result) public payable {
        address addr = msg.sender;
        uint256[] memory records = new uint256[](3);
        records[0] = matchId;
        records[1] = msg.value;
        records[2] = result;
        addressMatchsMap[addr][matchId] = records;
        matchAddressMap[matchId].push(addr);
        matchSums[matchId][result] +=  msg.value;
        uint256 total_sum = matchSums[matchId][0] + matchSums[matchId][1] + matchSums[matchId][2];
        matchOdds[matchId][0] = total_sum/matchSums[matchId][0];
        matchOdds[matchId][1] = total_sum/matchSums[matchId][1];
        matchOdds[matchId][2] = total_sum/matchSums[matchId][2];
    }
    
    //get match all guess address
    function getGuesserAddressByMatchId(uint256 matchId)  
    public  
    returns (address[] addrs)
    {
        addrs = matchAddressMap[matchId];
    }
    
    //get match all guess address
    function getMatchAddressInfo(uint256, address addr)
    public 
    returns (uint256, uint256, uint256)
    {
        uint256 matchId = addressMatchsMap[addr][matchId][0];
        uint256 money   = addressMatchsMap[addr][matchId][1];
        uint256 result  = addressMatchsMap[addr][matchId][2];
        return (matchId, money, result);
    }
    
    // Lottery match
    function lotteryMatch(address contractAddr,uint256 matchId) public {
        //mathch is over
        FifaData data = FifaData(contractAddr);
        uint256 match_state;
        uint256 home_team_goals;
        uint256 guest_team_goals;
        uint256 win_state = 0;
        (match_state, home_team_goals, guest_team_goals) = data.getMatchStateAndGoals(matchId);
        require(match_state == 2, "match is not over, try it later.");
        if(home_team_goals > guest_team_goals) {
            win_state = 0;
        } else if(home_team_goals == guest_team_goals) {
            win_state = 1;
        } else {
            win_state = 2;
        }
        address[] memory addrs = matchAddressMap[matchId];
        for(uint i = 0; i < addrs.length; i++) {
            address addr = addrs[i];
            uint256 money   = addressMatchsMap[addr][matchId][1];
            uint256 result  = addressMatchsMap[addr][matchId][2];
            if(win_state == result) {
                addr.transfer(money**matchOdds[matchId][win_state]);
            }
        }
    }
    
}
