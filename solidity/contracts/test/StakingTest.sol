// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.10;

/* solhint-disable no-global-import */
import "../interfaces/IStaking.sol";

/* solhint-disable custom-errors */

contract StakingTest is IStaking {
    mapping(string => uint256) public validatorShares;
    address public constant STAKING_ADDRESS =
        address(0x0000000000000000000000000000000000001003);

    function delegateV2(
        string memory _val,
        uint256 _amount
    ) external override returns (bool _result) {
        require(address(this).balance >= _amount, "insufficient balance");
        return IStaking(STAKING_ADDRESS).delegateV2(_val, _amount);
    }

    function undelegateV2(
        string memory _val,
        uint256 _amount
    ) external override returns (bool _result) {
        return IStaking(STAKING_ADDRESS).undelegateV2(_val, _amount);
    }

    function redelegateV2(
        string memory _valSrc,
        string memory _valDst,
        uint256 _amount
    ) external override returns (bool _result) {
        return
            IStaking(STAKING_ADDRESS).redelegateV2(_valSrc, _valDst, _amount);
    }

    function withdraw(string memory _val) external override returns (uint256) {
        return IStaking(STAKING_ADDRESS).withdraw(_val);
    }

    function approveShares(
        string memory _val,
        address _spender,
        uint256 _shares
    ) external override returns (bool) {
        return IStaking(STAKING_ADDRESS).approveShares(_val, _spender, _shares);
    }

    function transferShares(
        string memory _val,
        address _to,
        uint256 _shares
    ) external override returns (uint256, uint256) {
        return IStaking(STAKING_ADDRESS).transferShares(_val, _to, _shares);
    }

    function transferFromShares(
        string memory _val,
        address _from,
        address _to,
        uint256 _shares
    ) external override returns (uint256, uint256) {
        return
            IStaking(STAKING_ADDRESS).transferFromShares(
                _val,
                _from,
                _to,
                _shares
            );
    }

    function delegation(
        string memory _val,
        address _del
    ) public view override returns (uint256, uint256) {
        return IStaking(STAKING_ADDRESS).delegation(_val, _del);
    }

    function delegationRewards(
        string memory _val,
        address _del
    ) public view override returns (uint256) {
        return IStaking(STAKING_ADDRESS).delegationRewards(_val, _del);
    }

    function allowanceShares(
        string memory _val,
        address _owner,
        address _spender
    ) public view override returns (uint256) {
        return
            IStaking(STAKING_ADDRESS).allowanceShares(_val, _owner, _spender);
    }

    function slashingInfo(
        string memory _val
    ) external view override returns (bool _jailed, uint256 _missed) {
        return IStaking(STAKING_ADDRESS).slashingInfo(_val);
    }

    function validatorList(
        IStaking.ValidatorSortBy _sortBy
    ) external view override returns (string[] memory) {
        return IStaking(STAKING_ADDRESS).validatorList(_sortBy);
    }
}
