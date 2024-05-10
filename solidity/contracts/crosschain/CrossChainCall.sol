// SPDX-License-Identifier: Apache-2.0

/* solhint-disable one-contract-per-file */
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

library Encode {
    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "crossChain(address,string,uint256,uint256,bytes32,string)",
                _token,
                _receipt,
                _amount,
                _fee,
                _target,
                _memo
            );
    }

    // Deprecated: for fip20 only
    function fip20CrossChain(
        address _sender,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "fip20CrossChain(address,string,uint256,uint256,bytes32,string)",
                _sender,
                _receipt,
                _amount,
                _fee,
                _target,
                _memo
            );
    }

    function cancelSendToExternal(
        string memory _chain,
        uint256 _txid
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "cancelSendToExternal(string,uin256)",
                _chain,
                _txid
            );
    }

    function increaseBridgeFee(
        string memory _chain,
        uint256 _txid,
        address _token,
        uint256 _fee
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "increaseBridgeFee(string,uin256,address,uint256)",
                _chain,
                _txid,
                _token,
                _fee
            );
    }

    function bridgeCoinAmount(
        address _token,
        bytes32 _target
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "bridgeCoinAmount(address,bytes32)",
                _token,
                _target
            );
    }

    function bridgeCall(
        string memory _dstChainId,
        address _receiver,
        address[] memory _tokens,
        uint256[] memory _amounts,
        address _to,
        bytes memory _data,
        uint256 _value,
        bytes memory _memo
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "bridgeCall(string,address,address[],uint256[],address,bytes,uint256,bytes)",
                _dstChainId,
                _receiver,
                _tokens,
                _amounts,
                _to,
                _data,
                _value,
                _memo
            );
    }
}

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
        address _receiver,
        address[] memory _tokens,
        uint256[] memory _amounts,
        address _to,
        bytes memory _data,
        uint256 _value,
        bytes memory _memo
    ) internal returns (bool) {
        (bool result, bytes memory data) = CROSS_CHAIN_ADDRESS.call{
            value: msg.value
        }(
            Encode.bridgeCall(
                _dstChainId,
                _receiver,
                _tokens,
                _amounts,
                _to,
                _data,
                _value,
                _memo
            )
        );
        Decode.ok(result, data, "bridge call failed");
        return Decode.bridgeCall(data);
    }
}
