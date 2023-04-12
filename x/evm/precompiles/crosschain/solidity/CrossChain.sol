// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./ICrossChain.sol";
import "./Encode.sol";
import "./Decode.sol";

contract CrossChain is ICrossChain {
    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external payable virtual override returns (bool _result) {
        // precompile logic
        return true;
    }

    // Deprecated only for FIP20 token cross chain
    function fip20CrossChain(
        address _sender,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external virtual override returns (bool _result) {
        // precompile logic
        return true;
    }

    function cancelSendToExternal(
        string memory _chain,
        uint256 _txID
    ) external virtual override returns (bool _result) {
        // precompile logic
        return true;
    }

    function increaseBridgeFee(
        string memory _chain,
        uint256 _txID,
        address _token,
        uint256 _fee
    ) external payable virtual override returns (bool _result) {
        // precompile logic
        return true;
    }
}
