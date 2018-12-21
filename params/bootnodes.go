// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://c5ad436b28d07235fce2e8a60f7547a805ecdeec3a618f6d1e50e4cfe94c908d327d4c4b198e9432d11886e8f39c299767221e91e466842277fa732a1c6121f7@35.177.196.113:30302",
	"enode://844344c534ce108c235491e3bbd41161ab77804931c83531364a60386b6562133758634a1d5536cad9b70a3ac1949ced1e766d4d0c3db86207d1670dcb276511@35.177.196.113:30304",
	"enode://4c44acab3f6d477a900f9e2bb0951ea046d4a39823f4cb3eb51b79bf6a6cfb6bae991346b5621edd8b595c508e7c78036184110dbd0711db97260344f9e06a20@35.177.196.113:30301",
	//"enode://a7f7b4c4f9f6d029baaff7e7517aec60c646287927ef6c455886f51f8e38d1210b1170d59f20902189006988f5f10d292e57d9706e01c04f30fb20b3b3cc5087@172.31.30.194:30303",
	//"enode://957abf55e739456b9e4e7b0be196dceb3908e2b189b6cd715e25c1167a4d4ef616928016822f44ad15c207ce24e3c5494991f4b3dbe295b01f1dcf88cfc99454@172.31.29.101:30303",
	//"enode://c9409b25d25d9759eba40d52035b29b9ad4e33881379db2bf92d959e50ce24c4cc5d9dcc68219c207bf2c88abe6747b1f68df12598265ab614d12dc191ba701d@172.31.28.243:30303",
	//"enode://3f5ea64d729ecf364ed034546079efce2494d47f09bbfe7c33d9d294e5517bb7fa74cd555f0e314bdbe58ef79d72f326facd1c116e6b280ebf556acc6b720f04@172.31.27.16:30303",
	//"enode://5a6196e32ca1662c8c294ad4c3326fd8eab6e62b4e0669709210969474d6747a204179e60b87d1a510fe46e58a02c9b709f407e82a601df1060d8f91334490be@172.31.18.143:30303",
	//"enode://bd9bc9c30e223562bc7d59c7fe03be2684d9869e4474cbd7cd1f14b79bbecb4bdc5ddc833c7b263a88972297522ac6245c5515072bb123ebfe220fde0d96c37c@172.31.27.55:30303",
	//"enode://4c44acab3f6d477a900f9e2bb0951ea046d4a39823f4cb3eb51b79bf6a6cfb6bae991346b5621edd8b595c508e7c78036184110dbd0711db97260344f9e06a20@172.31.27.227:30301",
	//"enode://7b4f8f4fad0ff6063cd2ef9ec2591abfea001ef431be2a9158dd88b48eb4483ab656fbfd56a9ff479a4dbc0fe3e4cc3d0e264842fe2fc02464670785dd02c90a@172.31.30.194:30301",
}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{
	"enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303", // IE
	"enode://343149e4feefa15d882d9fe4ac7d88f885bd05ebb735e547f12e12080a9fa07c8014ca6fd7f373123488102fe5e34111f8509cf0b7de3f5b44339c9f25e87cb8@52.3.158.184:30303",  // INFURA
	"enode://b6b28890b006743680c52e64e0d16db57f28124885595fa03a562be1d2bf0f3a1da297d56b13da25fb992888fd556d4c1a27b1f39d531bde7de1921c90061cc6@159.89.28.211:30303", // AKASHA
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://06051a5573c81934c9554ef2898eb13b33a34b94cf36b202b69fde139ca17a85051979867720d4bdae4323d4943ddf9aeeb6643633aa656e0be843659795007a@35.177.226.168:30303",
	"enode://0cc5f5ffb5d9098c8b8c62325f3797f56509bff942704687b6530992ac706e2cb946b90a34f1f19548cd3c7baccbcaea354531e5983c7d1bc0dee16ce4b6440b@40.118.3.223:30304",
	"enode://1c7a64d76c0334b0418c004af2f67c50e36a3be60b5e4790bdac0439d21603469a85fad36f2473c9a80eb043ae60936df905fa28f1ff614c3e5dc34f15dcd2dc@40.118.3.223:30306",
	"enode://85c85d7143ae8bb96924f2b54f1b3e70d8c4d367af305325d30a61385a432f247d2c75c45c6b4a60335060d072d7f5b35dd1d4c45f76941f62a4f83b6e75daaf@40.118.3.223:30307",
}
