// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IBridgeCallback {
    function bridgeCallback(
        address,
        address,
        address[] memory,
        uint256[] memory,
        address,
        bytes memory,
        bytes memory,
        uint256
    ) external;

    function bridgeCallbackV1(
        address,
        address,
        address,
        address[] memory,
        uint256[] memory,
        bytes memory,
        uint256,
        uint256,
        uint256
    ) external;
}
