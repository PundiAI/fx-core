// SPDX-License-Identifier: Apache-2.0
/* solhint-disable one-contract-per-file */
pragma solidity ^0.8.0;

import {IBridgeCall} from "./IBridgeCall.sol";
import {FxBridgeBase} from "../bridge/FxBridgeBase.sol";

interface IFxBridgeLogic is IBridgeCall, FxBridgeBase {
    /* solhint-disable func-name-mixedcase */
    function state_fxBridgeId() external view returns (bytes32);
    function state_powerThreshold() external view returns (uint256);

    function state_lastEventNonce() external view returns (uint256);
    function state_lastOracleSetCheckpoint() external view returns (bytes32);
    function state_lastOracleSetNonce() external view returns (uint256);
    function state_lastBatchNonces(
        address _erc20Address
    ) external view returns (uint256);

    function bridgeTokens(uint256 _index) external view returns (address);
    function getTokenStatus(
        address _tokenAddr
    ) external view returns (TokenStatus memory);
    function version() external view returns (string memory);
    function state_lastBridgeCallNonces(
        uint256 _index
    ) external view returns (bool);
    /* ============== BSC FUNCTIONS =============== */
    function convert_decimals(
        address _erc20Address
    ) external view returns (uint8);
    /* solhint-disable func-name-mixedcase */

    function lastBatchNonce(
        address _erc20Address
    ) external view returns (uint256);

    function checkAssetStatus(address _tokenAddr) external view returns (bool);

    function addBridgeToken(
        address _tokenAddr,
        bytes32 _channelIBC,
        bool _isOriginated
    ) external returns (bool);

    function pauseBridgeToken(address _tokenAddr) external returns (bool);

    function activeBridgeToken(address _tokenAddr) external returns (bool);

    function updateOracleSet(
        address[] memory _newOracles,
        uint256[] memory _newPowers,
        uint256 _newOracleSetNonce,
        address[] memory _currentOracles,
        uint256[] memory _currentPowers,
        uint256 _currentOracleSetNonce,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s
    ) external;

    function sendToFx(
        address _tokenContract,
        bytes32 _destination,
        bytes32 _targetIBC,
        uint256 _amount
    ) external;

    function submitBatch(
        address[] memory _currentOracles,
        uint256[] memory _currentPowers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        uint256[] memory _amounts,
        address[] memory _destinations,
        uint256[] memory _fees,
        uint256[2] memory _nonceArray,
        address _tokenContract,
        uint256 _batchTimeout,
        address _feeReceive
    ) external;

    function submitBridgeCall(
        OracleSignatures calldata _curOracleSigns,
        uint256[2] calldata _nonceArray,
        BridgeCallData calldata _input
    ) external;

    function transferOwner(
        address _token,
        address _newOwner
    ) external returns (bool);

    /* ============== HELP FUNCTIONS =============== */

    function pause() external;

    function unpause() external;

    function getBridgeTokenList() external view returns (BridgeToken[] memory);

    /* =============== CHECKPOINTS =============== */

    function oracleSetCheckpoint(
        bytes32 _fxbridgeId,
        bytes32 _methodName,
        uint256 _oracleSetNonce,
        address[] memory _oracles,
        uint256[] memory _powers
    ) external pure returns (bytes32);

    function submitBatchCheckpoint(
        bytes32 _fxbridgeId,
        bytes32 _methodName,
        uint256[] memory _amounts,
        address[] memory _destinations,
        uint256[] memory _fees,
        uint256 _batchNonce,
        address _tokenContract,
        uint256 _batchTimeout,
        address _feeReceive
    ) external pure returns (bytes32);

    function bridgeCallCheckpoint(
        bytes32 _fxbridgeId,
        bytes32 _methodName,
        uint256 _nonce,
        BridgeCallData calldata _input
    ) external pure returns (bytes32);
}
