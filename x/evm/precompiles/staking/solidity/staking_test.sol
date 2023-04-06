// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

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
        (uint256 amount, uint256 reward, uint256 completionTime) = _undelegate(
            _val,
            _shares
        );
        validatorShares[_val] -= _shares;
        return (amount, reward, completionTime);
    }

    function withdraw(
        string memory _val
    ) external override returns (uint256) {
        uint256 amount = _withdraw(_val);
        return amount;
    }

    function approve(
        string memory _val,
        address _spender,
        uint256 _shares
    ) external override returns (bool) {
        bool success = _approve(_val, _spender, _shares);
        return success;
    }

    function transfer(
        string memory _val,
        address _to,
        uint256 _shares
    ) external override returns (uint256, uint256) {
        (uint256 token, uint256 reward) = _transfer(_val, _to, _shares);
        return (token, reward);
    }

    function transferFrom(
        string memory _val,
        address _from,
        address _to,
        uint256 _shares
    ) external override returns (uint256, uint256) {
        (uint256 token, uint256 reward) = _transferFrom(_val, _from, _to, _shares);
        return (token, reward);
    }

    function delegation(
        string memory _val,
        address _del
    ) public view override returns (uint256, uint256) {
        return _delegation(_val, _del);
    }

    function delegationRewards(
        string memory _val,
        address _del
    ) public view override returns (uint256) {
        return _delegationRewards(_val, _del);
    }

    function allowance(
        string memory _val,
        address _owner,
        address _spender
    ) public view override returns (uint256) {
        return _allowance(_val, _owner, _spender);
    }

}
