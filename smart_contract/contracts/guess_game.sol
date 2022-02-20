// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
contract GuessGame {
    uint256 num;
    address public owner;
    address public winer;

    constructor() {
        owner = msg.sender;
    }

    function SetNum(uint256 n) public{
        require(msg.sender == owner);
        if(0 < n && n < 100 && winer != owner){
            num = n;
            winer = owner;
        }
    }

    function GuessNum(uint256 n) public {
        require(msg.sender != owner);
        if (winer == owner){
            if (n == num){
                winer = msg.sender;
            }
        }

    }

}