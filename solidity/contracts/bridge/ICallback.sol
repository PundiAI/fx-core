// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

/* solhint-disable one-contract-per-file */
interface IBridgeCallback {
    function bridgeCallback(
        address,
        address,
        address[] memory,
        uint256[] memory,
        bytes memory,
        bytes memory
    ) external;
}

/* solhint-disable one-contract-per-file */
interface IRefundCallback {
    function refundCallback(
        uint256,
        address[] memory,
        uint256[] memory
    ) external;
}
