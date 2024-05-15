// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IRefundCallback {
    function refundCallback(
        uint256,
        address[] memory,
        uint256[] memory
    ) external;
}
