// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.12;

import "./Staking.sol";

contract staking_test is Staking {
    mapping(string => uint256) public validatorShares;

    function delegate(
        string memory _val
    ) external payable override returns (uint256, uint256) {
        (uint256 newShares, uint256 reward) = _delegate(_val);
        validatorShares[_val] += newShares;

        return (newShares, reward);
    }

    function undelegate(
        string memory _val,
        uint256 _shares
    ) external override returns (uint256, uint256, uint256) {
        (uint256 amount, uint256 reward, uint256 endTime) = _undelegate(
            _val,
            _shares
        );
        validatorShares[_val] -= _shares;
        return (amount, reward, endTime);
    }
}
