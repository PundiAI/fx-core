// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

library Decode {
    function crossChain(bytes memory data) internal pure returns (bool) {
        bool result = abi.decode(data, (bool));
        return result;
    }

    // Deprecated: for fip20 only
    function fip20CrossChain(bytes memory data) internal pure returns (bool) {
        bool result = abi.decode(data, (bool));
        return result;
    }

    function cancelSendToExternal(
        bytes memory data
    ) internal pure returns (bool) {
        bool result = abi.decode(data, (bool));
        return result;
    }

    function increaseBridgeFee(bytes memory data) internal pure returns (bool) {
        bool result = abi.decode(data, (bool));
        return result;
    }

    function bridgeCoinAmount(
        bytes memory data
    ) internal pure returns (uint256) {
        uint256 amount = abi.decode(data, (uint256));
        return amount;
    }

    function bridgeCall(bytes memory data) internal pure returns (bool) {
        return abi.decode(data, (bool));
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
