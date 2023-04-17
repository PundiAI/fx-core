// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

import "./Encode.sol";
import "./Decode.sol";

library CrossChainCall {
    address public constant CrossChainAddress = address(0x0000000000000000000000000000000000001004);

    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) internal returns (bool) {
        (bool result, bytes memory data) = CrossChainAddress.call{
                value: msg.value
            }(Encode.crossChain(_token, _receipt, _amount, _fee, _target, _memo));
        Decode.ok(result, data, "cross-chain failed");
        return Decode.crossChain(data);
    }

    function cancelSendToExternal(
        string memory _chain,
        uint256 _txID
    ) internal returns (bool) {
        (bool result, bytes memory data) = CrossChainAddress.call(
            Encode.cancelSendToExternal(_chain, _txID)
        );
        Decode.ok(result, data, "cancel send to external failed");
        return Decode.cancelSendToExternal(data);
    }

    function increaseBridgeFee(
        string memory _chain,
        uint256 _txID,
        address _token,
        uint256 _fee
    ) internal returns (bool) {
        (bool result, bytes memory data) = CrossChainAddress.call(
            Encode.increaseBridgeFee(_chain, _txID, _token, _fee)
        );
        Decode.ok(result, data, "increase bridge fee failed");
        return Decode.increaseBridgeFee(data);
    }
}
