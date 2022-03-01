// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./FIP20.sol";

contract WFX is FIP20 {

    constructor(string memory name_, string memory symbol_, uint8 decimals_) FIP20(name_, symbol_, decimals_){}

    receive() external payable {
        deposit();
    }

    fallback() external payable {
        deposit();
    }

    function deposit() public payable {
        balanceOf[msg.sender] += msg.value;
        emit Deposit(msg.sender, msg.value);
    }

    function withdraw(address payable to, uint256 value) public {
        require(balanceOf[msg.sender] >= value);
        balanceOf[msg.sender] -= value;
        to.transfer(value);
        emit Withdraw(msg.sender, to, value);
    }

    event Deposit(address indexed from, uint256 value);
    event Withdraw(address indexed from, address indexed to, uint256 value);
}