// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IBridgeCallback {
    function bridgeCallback(
        address _sender,
        address _receiver,
        address[] memory _tokens,
        uint256[] memory _amounts,
        bytes memory _data,
        bytes memory _memo
    ) external;
}
