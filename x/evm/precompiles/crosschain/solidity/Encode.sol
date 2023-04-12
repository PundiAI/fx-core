// SPDX-License-Identifier: UNLICENSED
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
}
