// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.12;

import "./Staking.sol";

contract staking_test is Staking {
    mapping(string => uint256) public validatorShares;

    function delegate(
        string memory _val,
        uint256 _amt
    ) external payable override returns (uint256) {
        uint256 newShares = _delegate(_val, _amt);
        validatorShares[_val] += newShares;

        return newShares;
    }

    function undelegate(
        string memory _val,
        uint256 _shares
    ) external override returns (uint256, uint) {
        (uint256 amount, uint256 endTime) = _undelegate(_val, _shares);
        validatorShares[_val] -= _shares;
        return (amount, endTime);
    }
}
