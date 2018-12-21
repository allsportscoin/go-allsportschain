pragma solidity ^0.4.25;

contract updatevote {

    address owner;

    modifier onlyOwner() {
        require(owner == msg.sender );
        _;
    }

    string private newVersion;
    string private curVersion;
    uint256 private finshNode;
    uint256 private mustPassNum;
    string private desc;
    bool private voteConfirm;
    address[] private condidate;
    mapping (address => string) private condidateState;

    constructor () public {
        owner = msg.sender;
        curVersion = "1.0.0";
        newVersion = curVersion;
        finshNode = 0;
        voteConfirm = false;
        mustPassNum = 0;
    }

    function putProposal (
        uint256 ifinshNode,
        string version,
        address[] icandidate,
        uint256 candidateLen,
        uint256 passNum,
        string  idesc
    )
        public
        onlyOwner
    {
        if (finshNode <= block.number) {
            if (voteConfirm) {
                curVersion = newVersion;
                voteConfirm = false;
            } else {
                newVersion = curVersion;
            }
            finshNode = 0;
            mustPassNum = 0;
            desc = "";
        }

        require (finshNode <= block.number
            && !strIsEqual(version, curVersion)
            && ifinshNode > block.number
            && passNum > 0
            && candidateLen >= passNum,
            "bad proposal");

        finshNode = ifinshNode;
        newVersion = version;
        condidate.length = 0;
        uint256 i = 0;
        for (i=0; i<candidateLen; i++) {
            condidate.push(icandidate[i]);
            condidateState[icandidate[i]] = curVersion;
        }
        mustPassNum = passNum;
        desc = idesc;
    }

    function cancalPoroposal() public onlyOwner {
        require (finshNode > block.number);
        newVersion = curVersion;
        voteConfirm = false;
        finshNode = 0;
        desc = "";
    }

    function voteUpdate(bool pass) public {
        require(!strIsEqual(newVersion, curVersion) && finshNode > block.number);
        bool    flag = false;
        uint256 passNum = 0;
        uint256 i;
        for (i=0; i<condidate.length; i++) {
            if (condidate[i] == msg.sender) {
                if (pass) {
                    condidateState[msg.sender] = newVersion;
                } else {
                    condidateState[msg.sender] = curVersion;
                }
                flag = true;
            }

            if (strIsEqual(condidateState[condidate[i]], newVersion)) {
                passNum++;
            }
        }

        require(flag);

        if (passNum >= mustPassNum) {
            voteConfirm = true;
        } else {
            voteConfirm = false;
        }
    }

    function getVersion () public view returns (string version) {
        if (!strIsEqual(curVersion, newVersion)
            && block.number >= finshNode
            && voteConfirm) {
            return newVersion;
        } else {
            return curVersion;
        }
    }

    function getProposalState()
        public
        view
        returns
    (
        uint256     ifinshNode,
        string      idesc,
        string      iversion,
        uint256     needPassNum,
        address[]   icondidate,
        bool[]      pass,
        uint256     len
    ) {
        if(finshNode <= block.number) {
            ifinshNode = 0;
            return;
        }

        icondidate = condidate;
        pass = new bool[](condidate.length);
        for (len=0; len<icondidate.length; len++) {
            pass[len] = strIsEqual(condidateState[icondidate[len]], newVersion);
        }
        ifinshNode = finshNode;
        needPassNum = mustPassNum;
        iversion = newVersion;
        idesc = desc;
    }

    function () public payable {
    }

    function strIsEqual(string str1, string str2) private pure returns(bool) {
        bytes32 str1Hash = keccak256(abi.encodePacked(str1));
        bytes32 str2Hash = keccak256(abi.encodePacked(str2));
        if(str1Hash == str2Hash) {
            return true;
        } else {
            return false;
        }
    }
}