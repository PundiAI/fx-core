// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IBridgeCall {
    function bridgeCall(
        string memory _dstChain,
        address _receiver,
        address[] memory _tokens,
        uint256[] memory _amounts,
        address _to,
        bytes memory _data,
        uint256 _value,
        bytes memory _memo
    ) external payable returns (uint256 _eventNonce);
}
