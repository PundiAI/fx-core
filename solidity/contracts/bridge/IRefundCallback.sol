// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IRefundCallback {
    function refundCallback(
        uint256 _eventNonce,
        address[] memory _tokens,
        uint256[] memory _amounts
    ) external;
}
