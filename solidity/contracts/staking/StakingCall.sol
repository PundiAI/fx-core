// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

import "./Encode.sol";
import "./Decode.sol";

library StakingCall {
    address public constant _stakingAddress = address(0x0000000000000000000000000000000000001003);

    function delegate(string memory _val, uint256 _value) internal returns (uint256, uint256) {
        (bool result, bytes memory data) = _stakingAddress.call{
                value: _value
            }(Encode.delegate(_val));
        Decode.ok(result, data, "delegate failed");
        return Decode.delegate(data);
    }

    function undelegate(
        string memory _val,
        uint256 _shares
    ) internal returns (uint256, uint256, uint256) {
        (bool result, bytes memory data) = _stakingAddress.call(
            Encode.undelegate(_val, _shares)
        );
        Decode.ok(result, data, "undelegate failed");
        return Decode.undelegate(data);
    }

    function withdraw(string memory _val) internal returns (uint256) {
        (bool result, bytes memory data) = _stakingAddress.call(
            Encode.withdraw(_val)
        );
        Decode.ok(result, data, "withdraw failed");
        return Decode.withdraw(data);
    }

    function approveShares(
        string memory _val,
        address _spender,
        uint256 _shares
    ) internal returns (bool) {
        (bool result, bytes memory data) = _stakingAddress.call(
            Encode.approveShares(_val, _spender, _shares)
        );
        Decode.ok(result, data, "approve shares failed");
        return Decode.approveShares(data);
    }

    function transferShares(
        string memory _val,
        address _to,
        uint256 _shares
    ) internal returns (uint256, uint256) {
        (bool result, bytes memory data) = _stakingAddress.call(
            Encode.transferShares(_val, _to, _shares)
        );
        Decode.ok(result, data, "transfer shares failed");
        return Decode.transferShares(data);
    }

    function transferFromShares(
        string memory _val,
        address _from,
        address _to,
        uint256 _shares
    ) internal returns (uint256, uint256) {
        (bool result, bytes memory data) = _stakingAddress.call(
            Encode.transferFromShares(_val, _from, _to, _shares)
        );
        Decode.ok(result, data, "transferFrom shares failed");
        return Decode.transferFromShares(data);
    }

    function delegation(
        string memory _val,
        address _del
    ) internal view returns (uint256, uint256) {
        (bool result, bytes memory data) = _stakingAddress.staticcall(
            Encode.delegation(_val, _del)
        );
        Decode.ok(result, data, "delegation failed");
        return Decode.delegation(data);
    }

    function delegationRewards(
        string memory _val,
        address _del
    ) internal view returns (uint256) {
        (bool result, bytes memory data) = _stakingAddress.staticcall(
            Encode.delegationRewards(_val, _del)
        );
        Decode.ok(result, data, "delegationRewards failed");
        return Decode.delegationRewards(data);
    }

    function allowanceShares(
        string memory _val,
        address _owner,
        address _spender
    ) internal view returns (uint256) {
        (bool result, bytes memory data) = _stakingAddress.staticcall(
            Encode.allowanceShares(_val, _owner, _spender)
        );
        Decode.ok(result, data, "allowance shares failed");
        return Decode.allowanceShares(data);
    }
}
