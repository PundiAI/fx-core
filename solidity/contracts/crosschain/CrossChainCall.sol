// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

/* solhint-disable no-global-import */
import "./Encode.sol";
import "./Decode.sol";

/* solhint-enable no-global-import */

library CrossChainCall {
    address public constant CROSS_CHAIN_ADDRESS =
        address(0x0000000000000000000000000000000000001004);

    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) internal returns (bool) {
        (bool result, bytes memory data) = CROSS_CHAIN_ADDRESS.call{
            value: msg.value
        }(Encode.crossChain(_token, _receipt, _amount, _fee, _target, _memo));
        Decode.ok(result, data, "cross-chain failed");
        return Decode.crossChain(data);
    }

    // Deprecated: for fip20 only
    function fip20CrossChain(
        address _sender,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) internal returns (bool) {
        // solhint-disable-next-line avoid-low-level-calls
        (bool result, bytes memory data) = CROSS_CHAIN_ADDRESS.call(
            Encode.fip20CrossChain(
                _sender,
                _receipt,
                _amount,
                _fee,
                _target,
                _memo
            )
        );
        Decode.ok(result, data, "fip-cross-chain failed");
        return Decode.fip20CrossChain(data);
    }

    function cancelSendToExternal(
        string memory _chain,
        uint256 _txID
    ) internal returns (bool) {
        // solhint-disable-next-line avoid-low-level-calls
        (bool result, bytes memory data) = CROSS_CHAIN_ADDRESS.call(
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
        // solhint-disable-next-line avoid-low-level-calls
        (bool result, bytes memory data) = CROSS_CHAIN_ADDRESS.call(
            Encode.increaseBridgeFee(_chain, _txID, _token, _fee)
        );
        Decode.ok(result, data, "increase bridge fee failed");
        return Decode.increaseBridgeFee(data);
    }

    function bridgeCoinAmount(
        address _token,
        bytes32 _target
    ) internal view returns (uint256) {
        (bool result, bytes memory data) = CROSS_CHAIN_ADDRESS.staticcall(
            Encode.bridgeCoinAmount(_token, _target)
        );
        Decode.ok(result, data, "bridge coin failed");
        return Decode.bridgeCoinAmount(data);
    }

    function bridgeCall(
        string memory _dstChainId,
        uint256 _gasLimit,
        address _receiver,
        address _to,
        bytes calldata _message,
        uint256 _value,
        bytes memory _asset
    ) internal returns (bool) {
        (bool result, bytes memory data) = CROSS_CHAIN_ADDRESS.call{
            value: msg.value
        }(
            Encode.bridgeCall(
                _dstChainId,
                _gasLimit,
                _receiver,
                _to,
                _message,
                _value,
                _asset
            )
        );
        Decode.ok(result, data, "bridge call failed");
        return Decode.bridgeCall(data);
    }
}
