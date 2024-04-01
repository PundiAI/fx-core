// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

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
        uint256 _gasLimit,
        address _receiver,
        address _to,
        bytes calldata _message,
        uint256 _value,
        bytes memory _asset
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "bridgeCall(string,uint256,address,address,bytes,uint256,bytes)",
                _dstChainId,
                _gasLimit,
                _receiver,
                _to,
                _message,
                _value,
                _asset
            );
    }
}
