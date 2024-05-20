// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

import {IBridgeCall} from "./IBridgeCall.sol";

interface ICrossChain is IBridgeCall {
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

    function bridgeCoinAmount(
        address _token,
        bytes32 _target
    ) external view returns (uint256 _amount);

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

    event BridgeCallEvent(
        address indexed _sender,
        address indexed _receiver,
        address indexed _to,
        address _txOrigin,
        uint256 _value,
        uint256 _eventNonce,
        string _dstChain,
        address[] _tokens,
        uint256[] _amounts,
        bytes _data,
        bytes _memo
    );
}
