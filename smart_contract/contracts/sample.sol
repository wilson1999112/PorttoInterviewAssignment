// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
contract Coin {
    address owner;
    mapping (address => uint) balances;
    constructor() {
        owner = msg.sender;
    }
    function mint(uint amount) public{
        require (msg.sender == owner);
        balances[msg.sender] += amount;
    }
    function send(address receiver, uint amount) public{
        require (balances[msg.sender] >= amount);
        balances[msg.sender] -= amount;
        balances[receiver] += amount;
    }
    function getBalance(address addr) public view returns (uint balance) {
        return balances[addr];
    }
}