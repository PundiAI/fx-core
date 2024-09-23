// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.10;

import {IBridgeCall} from "./IBridgeCall.sol";

// NOTE: if using an interface to invoke the precompiled contract
// need to use solidity version 0.8.10 and later.
interface ICrossChain is IBridgeCall {
    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external payable returns (bool _result);

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

    function cancelPendingBridgeCall(
        string memory _chain,
        uint256 _txID
    ) external returns (bool _result);

    function addPendingPoolRewards(
        string memory _chain,
        uint256 _txID,
        address _token,
        uint256 _reward
    ) external payable returns (bool _result);

    function bridgeCoinAmount(
        address _token,
        bytes32 _target
    ) external view returns (uint256 _amount);

    function executeClaim(
        string memory _chain,
        uint256 _eventNonce
    ) external returns (bool _result);

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

    event CancelPendingBridgeCallEvent(
        address indexed _sender,
        string chain,
        uint256 txID
    );

    event AddPendingPoolRewardsEvent(
        address indexed sender,
        address indexed token,
        string chain,
        uint256 txID,
        uint256 reward
    );

    event ExecuteClaimEvent(
        address indexed _sender,
        uint256 _eventNonce,
        string _chain
    );
}
