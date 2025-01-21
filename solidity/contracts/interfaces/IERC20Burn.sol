// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

interface IERC20Burn {
    function burn(uint256 amount) external;
    function burnFrom(address account, uint256 amount) external;
}
