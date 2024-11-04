// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

import {IBridgeCall} from "./IBridgeCall.sol";

interface IFxBridgeLogic is IBridgeCall {
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
    function tokenStatus(
        address _tokenAddr
    ) external view returns (TokenStatus memory);
    function version() external view returns (string memory);
    function state_lastBridgeCallNonces(
        uint256 _index
    ) external view returns (bool);
    /* solhint-disable func-name-mixedcase */

    struct TokenStatus {
        bool isOriginated;
        bool isActive;
        bool isExist;
    }

    struct TransferInfo {
        uint256 amount;
        address destination;
        uint256 fee;
        address exchange;
        uint256 minExchange;
    }

    struct BridgeToken {
        address addr;
        string name;
        string symbol;
        uint8 decimals;
    }

    struct BridgeCallData {
        address sender;
        address refund;
        address[] tokens;
        uint256[] amounts;
        address to;
        bytes data;
        bytes memo;
        uint256 timeout;
        uint256 eventNonce;
    }

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
        address[] memory _currentOracles,
        uint256[] memory _currentPowers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        uint256[2] memory _nonceArray,
        BridgeCallData memory _input
    ) external;

    function transferOwner(
        address _token,
        address _newOwner
    ) external returns (bool);

    /* =============== QUERY FUNCTIONS =============== */

    function lastBatchNonce(
        address _erc20Address
    ) external view returns (uint256);

    function checkAssetStatus(address _tokenAddr) external view returns (bool);

    /* ============== HELP FUNCTIONS =============== */

    function makeCheckpoint(
        address[] memory _oracles,
        uint256[] memory _powers,
        uint256 _oracleSetNonce,
        bytes32 _fxBridgeId
    ) external pure returns (bytes32);

    function checkOracleSignatures(
        address[] memory _currentOracles,
        uint256[] memory _currentPowers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        bytes32 _theHash,
        uint256 _powerThreshold
    ) external pure;

    function pause() external;

    function unpause() external;

    function getBridgeTokenList() external view returns (BridgeToken[] memory);

    /* ============== BSC FUNCTIONS =============== */
    function convert_decimals(
        address _erc20Address
    ) external view returns (uint8);

    /* =============== CHECKPOINTS =============== */

    function oracleSetCheckpoint(
        bytes32 _fxbridgeId,
        bytes32 _methodName,
        uint256 _oracleSetNonce,
        address[] memory _oracles,
        uint256[] memory _powers
    ) external returns (bytes32);

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
    ) external returns (bytes32);

    function bridgeCallCheckpoint(
        bytes32 _fxbridgeId,
        bytes32 _methodName,
        address _sender,
        address _refund,
        address[] memory _tokens,
        uint256[] memory _amounts,
        address _to,
        bytes memory _data,
        bytes memory _memo,
        uint256 _nonce,
        uint256 _timeout,
        uint256 _eventNonce
    ) external returns (bytes32);

    /* =============== EVENTS =============== */

    event TransactionBatchExecutedEvent(
        uint256 indexed _batchNonce,
        address indexed _token,
        uint256 _eventNonce
    );
    event SendToFxEvent(
        address indexed _tokenContract,
        address indexed _sender,
        bytes32 indexed _destination,
        bytes32 _targetIBC,
        uint256 _amount,
        uint256 _eventNonce
    );
    event AddBridgeTokenEvent(
        address indexed _tokenContract,
        string _name,
        string _symbol,
        uint8 _decimals,
        uint256 _eventNonce,
        bytes32 _channelIBC
    );
    event OracleSetUpdatedEvent(
        uint256 indexed _newOracleSetNonce,
        uint256 _eventNonce,
        address[] _oracles,
        uint256[] _powers
    );
    event TransferOwnerEvent(address _token, address _newOwner);

    event BridgeCallEvent(
        address indexed _sender,
        address indexed _refund,
        address indexed _to,
        address _txOrigin,
        uint256 _eventNonce,
        string _dstChain,
        address[] _tokens,
        uint256[] _amounts,
        bytes _data,
        uint256 _quoteId,
        bytes _memo
    );

    event SubmitBridgeCallEvent(
        address indexed _txOrigin,
        uint256 _nonce,
        uint256 _eventNonce,
        bool _success,
        bytes _cause
    );
}
