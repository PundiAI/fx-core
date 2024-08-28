// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.10;

/* solhint-disable no-global-import */
import "../staking/IStaking.sol";
import "../staking/StakingCall.sol";

/* solhint-enable no-global-import */

/* solhint-disable custom-errors */

contract StakingTest is IStaking {
    mapping(string => uint256) public validatorShares;

    function delegate(
        string memory _val
    ) external payable override returns (uint256, uint256) {
        (uint256 newShares, uint256 reward) = StakingCall.delegate(
            _val,
            msg.value
        );
        validatorShares[_val] += newShares;
        return (newShares, reward);
    }

    function delegateV2(
        string memory _val,
        uint256 _amount
    ) external payable override returns (bool _result) {
        require(address(this).balance >= _amount, "insufficient balance");
        return IStaking(StakingCall.STAKING_ADDRESS).delegateV2(_val, _amount);
    }

    function undelegate(
        string memory _val,
        uint256 _shares
    ) external override returns (uint256, uint256, uint256) {
        (uint256 amount, uint256 reward, uint256 completionTime) = StakingCall
            .undelegate(_val, _shares);
        validatorShares[_val] -= _shares;
        return (amount, reward, completionTime);
    }

    function undelegateV2(
        string memory _val,
        uint256 _amount
    ) external override returns (bool _result) {
        return
            IStaking(StakingCall.STAKING_ADDRESS).undelegateV2(_val, _amount);
    }

    function redelegate(
        string memory _valSrc,
        string memory _valDst,
        uint256 _shares
    ) external override returns (uint256, uint256, uint256) {
        (uint256 amount, uint256 reward, uint256 completionTime) = StakingCall
            .redelegate(_valSrc, _valDst, _shares);
        validatorShares[_valSrc] -= _shares;
        validatorShares[_valDst] += _shares;
        return (amount, reward, completionTime);
    }

    function redelegateV2(
        string memory _valSrc,
        string memory _valDst,
        uint256 _amount
    ) external override returns (bool _result) {
        return
            IStaking(StakingCall.STAKING_ADDRESS).redelegateV2(
                _valSrc,
                _valDst,
                _amount
            );
    }

    function withdraw(string memory _val) external override returns (uint256) {
        uint256 amount = StakingCall.withdraw(_val);
        return amount;
    }

    function approveShares(
        string memory _val,
        address _spender,
        uint256 _shares
    ) external override returns (bool) {
        bool success = StakingCall.approveShares(_val, _spender, _shares);
        return success;
    }

    function transferShares(
        string memory _val,
        address _to,
        uint256 _shares
    ) external override returns (uint256, uint256) {
        (uint256 token, uint256 reward) = StakingCall.transferShares(
            _val,
            _to,
            _shares
        );
        return (token, reward);
    }

    function transferFromShares(
        string memory _val,
        address _from,
        address _to,
        uint256 _shares
    ) external override returns (uint256, uint256) {
        (uint256 token, uint256 reward) = StakingCall.transferFromShares(
            _val,
            _from,
            _to,
            _shares
        );
        return (token, reward);
    }

    function delegation(
        string memory _val,
        address _del
    ) public view override returns (uint256, uint256) {
        return StakingCall.delegation(_val, _del);
    }

    function delegationRewards(
        string memory _val,
        address _del
    ) public view override returns (uint256) {
        return StakingCall.delegationRewards(_val, _del);
    }

    function allowanceShares(
        string memory _val,
        address _owner,
        address _spender
    ) public view override returns (uint256) {
        return StakingCall.allowanceShares(_val, _owner, _spender);
    }

    function slashingInfo(
        string memory _val
    ) external view override returns (bool _jailed, uint256 _missed) {
        return IStaking(StakingCall.STAKING_ADDRESS).slashingInfo(_val);
    }

    function validatorList(
        IStaking.ValidatorSortBy _sortBy
    ) external view override returns (string[] memory) {
        return IStaking(StakingCall.STAKING_ADDRESS).validatorList(_sortBy);
    }
}
