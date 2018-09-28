pragma solidity ^0.4.19;

contract MatchData {
    
    address internal owner;

    modifier onlyOwner() {
        require(owner == msg.sender );
        _;
    }
    
    modifier idVaild(uint256 id) {
        require(id > 0 && id < matchs.length);
        _;
    }
    
    //data
    //string name;
    
    struct Match {
        string name;
        uint256 price;
        string desc;
        uint256 state;
        uint256 startTime;
        uint256 endTime;
        string  homeName;
        string  guestName;
        uint8   homeScore;
        uint8   guestScore;
    }
    
    Match[] private matchs;
    mapping (string => uint256) matchIndex; 
    
    //function
    constructor () payable public{
        owner = msg.sender;
        matchs.push(Match("", 0, "", 0, 0, 0, "", "", 0, 0));
    }
    
    function AddMatch(
        string name, 
        uint256 price, 
        string desc,
        uint256 state,
        uint256 startTime,
        uint256 endTime,
        string  homeName,
        string  guestName
    ) public onlyOwner returns (uint256) {
        require(matchIndex[name] == 0);//name is no-repeat
        matchs.push(Match(name,price,desc,state,startTime,endTime,homeName,guestName,0,0));
        uint256 index = matchs.length - 1;
        matchIndex[name] = index;
        return matchs.length - 1;
    }
    
    function GetMatchMaxId() public view returns (uint256) {
        return matchs.length - 1;
    }
    
    function GetMatchInfo(
        uint256 id
    ) public view idVaild(id) 
    returns (
        string name, 
        //uint256 price, 
        //string desc, 
        //uint256 state, 
        uint256 startTime, 
        uint256 endTime,
        string  homeName,
        string  guestName) {
            
        name = matchs[id].name;
        //price = matchs[id].price;
        //desc = matchs[id].desc;
        //state = matchs[id].state;
        startTime = matchs[id].startTime;
        endTime = matchs[id].endTime;
        homeName = matchs[id].homeName;
        guestName = matchs[id].guestName;
    }
    
    function GetMatchInfo(
        string findName
    ) public view 
    returns (
        string, 
        //uint256, 
        //string, 
        //uint256, 
        uint256, 
        uint256,
        string, 
        string) {
            
        uint256 id = matchIndex[findName];
        return GetMatchInfo(id);
    }
    
    function GetMatchState(
        uint256 id
    ) public view idVaild(id) returns (uint256) {
        return matchs[id].state;
    }
    
    function GetMatchState(
        string name
    ) public view returns (uint256) {
        uint256 id = matchIndex[name];
        return GetMatchState(id);
    }
    
    function GetMatchPrice(
        uint256 id
    ) public view idVaild(id) returns (uint256) {
        return matchs[id].price;
    }
    
    function GetMatchPrice(
        string findName
    ) public view returns (uint256) {
        uint256 id = matchIndex[findName];
        return GetMatchPrice(id);
    }
    
    function GetMatchDesc(
        uint256 id
    ) public view idVaild(id) returns (string) {
        return matchs[id].desc;
    }
    
    function GetMatchDesc(
        string findName
    ) public view returns (string) {
        uint256 id = matchIndex[findName];
        return GetMatchDesc(id);
    }
    
    function GetMatchTime(
        uint256 id
    ) public view idVaild(id) returns (uint256 startTime, uint256 endTime) {
        startTime = matchs[id].startTime;
        endTime = matchs[id].endTime;
    }
    
    function GetMatchTime(
        string findName
    ) public view returns (uint256, uint256) {
        uint256 id = matchIndex[findName];
        return GetMatchTime(id);
    }
    
    function GetMatchScroe(
        uint256 id
    ) public payable idVaild(id)  returns  (uint256 state, uint8 homeScore, uint8 guestScore) {
        require(id < matchs.length);
        require(msg.value >= matchs[id].price);
        
        state = matchs[id].state;
        homeScore = matchs[id].homeScore;
        guestScore = matchs[id].guestScore;
        
        if (msg.value > matchs[id].price) {
            address(msg.sender).transfer(msg.value - matchs[id].price);
        }
    }
    
    /*
    function GetMatchScroe(
        string findName
    ) public payable returns  (uint256 , uint8 , uint8 ) {
        uint256 id = matchIndex[findName];
        return GetMatchScroe(id);
    }
    */
    
    function GetMatchInfoByOwner(
        uint256 id
    ) public view onlyOwner idVaild(id) 
    returns (
        string name, 
        uint256 price, 
        string desc, 
        uint256 state,
        string homeName,
        string guestName,
        uint8 homeScore,
        uint8 guestScore
    ){
        name = matchs[id].name;
        price = matchs[id].price;
        desc = matchs[id].desc;
        state = matchs[id].state;
        homeScore = matchs[id].homeScore;
        homeName  = matchs[id].homeName;
        guestScore = matchs[id].guestScore;
        guestName = matchs[id].guestName;
    }
    
    function SetMatchInfo(
        uint256 id, 
        uint256 price, 
        string desc, 
        uint256 state) public onlyOwner idVaild(id)
    {
        matchs[id].price = price;
        matchs[id].desc = desc;
        matchs[id].state = state;
    }
    
    function SetMatchInfo(
        string findName, 
        uint256 price, 
        string desc, 
        uint256 state) public onlyOwner
    {
        uint256 id = matchIndex[findName];
        SetMatchInfo(id,price,desc,state);
    }
    
    function SetMatchState(
        uint256 id, 
        uint256 state
    ) public onlyOwner idVaild(id) 
    {
        matchs[id].state = state;
    }
    
    function SetMatchState(
        string findName,
        uint256 state
    ) public onlyOwner
    {
        uint256 id = matchIndex[findName];
        SetMatchState(id, state);
    }
    
    function SetMatchScore(
        uint256 id, 
        uint8 homeScore,
        uint8 guestScore) public onlyOwner idVaild(id) {
            
        matchs[id].homeScore = homeScore;
        matchs[id].guestScore = guestScore;
    }
    
    function SetMatchScore(
        string findName,
        uint8 homeScore,
        uint8 guestScore) public onlyOwner {
            
        uint256 id = matchIndex[findName];
        SetMatchScore(id, homeScore, guestScore);
    }
    
    function SetMatchTime(
        uint256 id, 
        uint256 startTime,
        uint256 endTime) public onlyOwner idVaild(id) {
            
        matchs[id].startTime = startTime;
        matchs[id].endTime = endTime;
    }
    
    function SetMatchTime(
        string findName,
        uint256 startTime,
        uint256 endTime) public onlyOwner{
            
        uint256 id = matchIndex[findName];
        SetMatchTime(id, startTime, endTime);
    }
    
    function SetMatchPrice(
        uint256 id, 
        uint256 price) public onlyOwner idVaild(id) {
        matchs[id].price = price;
    }
    
    function SetMatchPrice(
        string findName,
        uint256 price) public onlyOwner {
            
        uint256 id = matchIndex[findName];
        SetMatchPrice(id, price);
    }
    
    function GetBalance() public view returns (uint256) {
        return address(this).balance;
    }

    function Draw(
        address drawer, 
        uint256 val
    ) public onlyOwner{
        if (val < address(this).balance) {
            drawer.transfer(val);
        } else {
            drawer.transfer(address(this).balance);
        }
    }
    
    function Draw(
        address drawer
    ) public onlyOwner{
        drawer.transfer(address(this).balance);
    }
    
    function () public payable {
    }
}

contract FootballGuess  {
    
    address owner;
    address contractAddr;
    
    struct Match {
        uint256 id;
        string official_name;
        string home_team;
        string guest_team;
        uint256 start_time;
        uint256 end_time;
        uint256 state;
        uint256 home_team_goals;
        uint256 guest_team_goals;
        
        uint256 home_sum;
        uint256 draw_sum;
        uint256 away_sum;
        
        uint256 odds_home;
        uint256 odds_draw;
        uint256 odds_away;
    }
    
    Match[] matchsIdMap;
    
    struct match_addr_guess_info {
        uint256 home_sum;
        uint256 draw_sum;
        uint256 away_sum;
    }
    
    mapping(uint256 => mapping(address => match_addr_guess_info)) matchAddrGuessInfo;
    mapping(uint256 => address[])  matchAddressMap;
    mapping(uint256 => mapping(address => bool)) matchAddrGuessInfoBool;
    
    constructor() public payable {
        owner = msg.sender;
        matchsIdMap.push(Match(0, "init", "t_1", "t_2", 1537607619, 1537607619, 0,0,0, 0,0,0, 1,1,1));
    }

    modifier onlyOwner() {
        require(
            msg.sender == owner,
            "Only owner can call this."
        );
        _;
    }
    
    //get data contract address
    function getDataContractAddr() public view returns(address) {
        return contractAddr;
    }
    //set data contract address
    function setDataContractAddr(address addr) public onlyOwner {
        contractAddr = addr;
    }
    
    function GetBalance() public view returns (uint256) {
        return address(this).balance;
    }
    
    function Pay() public payable{
    }
    
    //refresh Match Info by owner
    function refreshMatchInfo() public payable returns (bool) {
        bool refresh = false;
        MatchData data = MatchData(contractAddr);
        
        uint256 max_id = data.GetMatchMaxId();
        string memory official_name = "";
        string memory home_team = "";
        string memory guest_team = "";
        uint256 start_time = 0;
        uint256 end_time = 0;
        uint256 len = matchsIdMap.length-1;
        if(max_id > len) {
            uint256 i = matchsIdMap.length;
            for(; i <= max_id; i++) {
                //store all info from data layer
                (official_name, start_time, end_time, home_team, guest_team) = data.GetMatchInfo(i);
                matchsIdMap.push(Match(i, official_name, home_team, guest_team, start_time, end_time, 0,0,0, 0,0,0, 1,1,1));
            }
            refresh = true;
        }
        return refresh;
    }
    
    //get match max id
    function getMatchMaxId() public view returns(uint256) {
        return matchsIdMap.length-1;
    }
    
    //get Match Info By Id From Local
    function getMatchInfoById(uint256 id)
    public view
    returns (uint256 match_id, string official_name, string home_team, string guest_team,
        uint256 start_time, uint256 end_time, uint256 state, uint256 home_team_goals, uint256 guest_team_goals,
        uint256 home_sum, uint256 draw_sum, uint256 away_sum
    ) 
    {
        require(
            id <= matchsIdMap.length,
            "id is too big now."
        );
        MatchData data = MatchData(contractAddr);
        state = data.GetMatchState(id);
        match_id = matchsIdMap[id].id;
        official_name = matchsIdMap[id].official_name; 
        home_team = matchsIdMap[id].home_team;
        guest_team = matchsIdMap[id].guest_team;
        start_time = matchsIdMap[id].start_time;
        end_time = matchsIdMap[id].end_time;
        home_team_goals = matchsIdMap[id].home_team_goals; 
        guest_team_goals = matchsIdMap[id].guest_team_goals;
        home_sum = matchsIdMap[id].home_sum;
        draw_sum = matchsIdMap[id].draw_sum;
        away_sum = matchsIdMap[id].away_sum;
    }
    
    //add match guess record
    function addMatchGuessRecord(uint256 matchId, uint256 result) public payable {
        MatchData data = MatchData(contractAddr);
        uint256 state = data.GetMatchState(matchId);
        require(state == 0, "match is start, can no guess!");
        
        address addr = msg.sender;
        if(matchAddrGuessInfoBool[matchId][addr]  == false){ 
            matchAddrGuessInfo[matchId][addr] = match_addr_guess_info(0,0,0);
            matchAddressMap[matchId].push(addr);
            matchAddrGuessInfoBool[matchId][addr] = true;
        } else {
           //alread exist, do not things
        }
        
        if(result == 0) {
            matchsIdMap[matchId].home_sum  += msg.value;
            matchAddrGuessInfo[matchId][addr].home_sum += msg.value;
        } else if (result == 1) {
            matchsIdMap[matchId].draw_sum  += msg.value;
            matchAddrGuessInfo[matchId][addr].draw_sum += msg.value;
        } else  {
            matchsIdMap[matchId].away_sum  += msg.value;
            matchAddrGuessInfo[matchId][addr].away_sum += msg.value;
        }
        
        uint256 total_sum = matchsIdMap[matchId].home_sum + matchsIdMap[matchId].draw_sum + matchsIdMap[matchId].away_sum;
        if (matchsIdMap[matchId].home_sum != 0) {
            matchsIdMap[matchId].odds_home = total_sum/matchsIdMap[matchId].home_sum;
        }
        if (matchsIdMap[matchId].draw_sum != 0) {
            matchsIdMap[matchId].odds_draw = total_sum/matchsIdMap[matchId].draw_sum;
        }
        if (matchsIdMap[matchId].away_sum != 0) {
            matchsIdMap[matchId].odds_away = total_sum/matchsIdMap[matchId].away_sum;
        }
    }
    
    //get match all guess address
    function getGuesserAddressByMatchId(uint256 matchId)  
    public view
    returns (address[] addrs)
    {
        addrs = matchAddressMap[matchId];
    }
    
    //get match address guess info
    function getMatchAddressRecordInfo(uint256 matchId, address addr)
    public view
    returns (uint256 home_sum, uint256 draw_sum, uint256 away_sum)
    {
         home_sum = matchAddrGuessInfo[matchId][addr].home_sum;
         draw_sum = matchAddrGuessInfo[matchId][addr].draw_sum;
         away_sum = matchAddrGuessInfo[matchId][addr].away_sum;
    }
    
    // Lottery match
    function lotteryMatch(uint256 matchId) 
    public payable
    returns(uint256 state, uint256 win_state)
    {

        //mathch is over
        MatchData data = MatchData(contractAddr);
        uint256 res_state;
        uint256 res_home_gloas;
        uint256 res_guest_gloas;
        (res_state , res_home_gloas, res_guest_gloas) = data.GetMatchScroe.value(10)(matchId);
        
        win_state = 0;
        state = res_state;
        require(res_state == 2, "match is not over, try it later.");
        require(matchsIdMap[matchId].state != 2, "match already draw a lottery");
        uint256 odds = 0;
        if(res_home_gloas > res_guest_gloas) {
            win_state = 0;
            odds = matchsIdMap[matchId].odds_home;
        } else if(res_home_gloas == res_guest_gloas) {
            win_state = 1;
            odds = matchsIdMap[matchId].odds_draw;
        } else if (res_home_gloas < res_guest_gloas) {
            win_state = 2;
            odds = matchsIdMap[matchId].odds_away;
        }
        
        address[] memory addrs = matchAddressMap[matchId];
        
        for(uint i = 0; i < addrs.length; i++) {
            address addr = addrs[i];
            uint256 money = 0;
            if(win_state == 0) {
                money = matchAddrGuessInfo[matchId][addr].home_sum;
                addr.transfer(money*odds);
            } else if (win_state == 1) {
                money = matchAddrGuessInfo[matchId][addr].draw_sum;
                addr.transfer(money*odds);
            } else if (win_state == 1){
                money = matchAddrGuessInfo[matchId][addr].away_sum;
                addr.transfer(money*odds);
            }
        }
        matchsIdMap[matchId].state = 2;
    }
}
