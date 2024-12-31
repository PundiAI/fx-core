// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

// NOTE: if using an interface to invoke the precompiled contract
// need to use solidity version 0.8.10 and later.
interface IBridgeCall {
    function bridgeCall(
        string memory _dstChain,
        address _refund,
        address[] memory _tokens,
        uint256[] memory _amounts,
        address _to,
        bytes memory _data,
        uint256 _quoteId,
        uint256 _gasLimit,
        bytes memory _memo
    ) external payable returns (uint256 _eventNonce);

    event BridgeCallEvent(
        address indexed _sender,
        address indexed _refund,
        address indexed _to,
        address _txOrigin,
        uint256 _eventNonce,
        string _dstChain,
        address[] _tokens,
        uint256[] _amounts,
        bytes _data,
        uint256 _quoteId,
        uint256 _gasLimit,
        bytes _memo
    );
}
