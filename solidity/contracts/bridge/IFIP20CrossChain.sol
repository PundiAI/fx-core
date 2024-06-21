// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

// NOTE: if using an interface to invoke the precompiled contract
// need to use solidity version 0.8.10 and later.
interface IFIP20CrossChain {
    // Deprecated: for fip20 only
    function fip20CrossChain(
        address _sender,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external returns (bool _result);
}
