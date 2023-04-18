// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

interface ICrossChain {
    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external payable returns (bool _result);

    // Deprecated: for fip20 only
    function fip20CrossChain(
        address _sender,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external returns (bool _result);

    function cancelSendToExternal(
        string memory _chain,
        uint256 _txID
    ) external returns (bool _result);

    function increaseBridgeFee(
        string memory _chain,
        uint256 _txID,
        address _token,
        uint256 _fee
    ) external payable returns (bool _result);

    event CrossChain(
        address indexed sender,
        address indexed token,
        string denom,
        string receipt,
        uint256 amount,
        uint256 fee,
        bytes32 target,
        string memo
    );

    event CancelSendToExternal(
        address indexed sender,
        string chain,
        uint256 txID
    );

    event IncreaseBridgeFee(
        address indexed sender,
        address indexed token,
        string chain,
        uint256 txID,
        uint256 fee
    );
}
