// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

library Decode {
    function delegate(
        bytes memory data
    ) internal pure returns (uint256, uint256) {
        (uint256 shares, uint256 reward) = abi.decode(data, (uint256, uint256));
        return (shares, reward);
    }

    function undelegate(
        bytes memory data
    ) internal pure returns (uint256, uint256, uint256) {
        (uint256 amount, uint256 reward, uint256 completionTime) = abi.decode(
            data,
            (uint256, uint256, uint256)
        );
        return (amount, reward, completionTime);
    }

    function redelegate(
        bytes memory data
    ) internal pure returns (uint256, uint256, uint256) {
        (uint256 amount, uint256 reward, uint256 completionTime) = abi.decode(
            data,
            (uint256, uint256, uint256)
        );
        return (amount, reward, completionTime);
    }

    function withdraw(bytes memory data) internal pure returns (uint256) {
        uint256 reward = abi.decode(data, (uint256));
        return reward;
    }

    function approveShares(bytes memory data) internal pure returns (bool) {
        bool result = abi.decode(data, (bool));
        return result;
    }

    function transferShares(
        bytes memory data
    ) internal pure returns (uint256, uint256) {
        (uint256 amount, uint256 reward) = abi.decode(data, (uint256, uint256));
        return (amount, reward);
    }

    function transferFromShares(
        bytes memory data
    ) internal pure returns (uint256, uint256) {
        (uint256 amount, uint256 reward) = abi.decode(data, (uint256, uint256));
        return (amount, reward);
    }

    function delegation(
        bytes memory data
    ) internal pure returns (uint256, uint256) {
        (uint256 delegateAmount, uint256 shares) = abi.decode(
            data,
            (uint256, uint256)
        );
        return (delegateAmount, shares);
    }

    function delegationRewards(
        bytes memory data
    ) internal pure returns (uint256) {
        uint256 delegateRewardsAmount = abi.decode(data, (uint256));
        return delegateRewardsAmount;
    }

    function allowanceShares(
        bytes memory data
    ) internal pure returns (uint256) {
        uint256 allowanceAmount = abi.decode(data, (uint256));
        return allowanceAmount;
    }

    function ok(
        bool _result,
        bytes memory _data,
        string memory _msg
    ) internal pure {
        if (!_result) {
            string memory errMsg = abi.decode(_data, (string));
            if (bytes(_msg).length < 1) {
                revert(errMsg);
            }
            revert(string(abi.encodePacked(_msg, ": ", errMsg)));
        }
    }
}
