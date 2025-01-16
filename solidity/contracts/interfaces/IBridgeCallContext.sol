// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IBridgeCallContext {
    function onBridgeCall(
        address _sender,
        address _refund,
        address[] memory _tokens,
        uint256[] memory _amounts,
        bytes memory _data,
        bytes memory _memo
    ) external;

    function onRevert(uint256 nonce, bytes memory _msg) external;
}
