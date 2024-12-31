// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

interface FxBridgeBase {
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

    struct OracleSignatures {
        address[] oracles;
        uint256[] powers;
        bytes32[] r;
        bytes32[] s;
        uint8[] v;
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
        uint256 gasLimit;
        uint256 eventNonce;
    }

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

    event SubmitBridgeCallEvent(
        address indexed _txOrigin,
        uint256 _nonce,
        uint256 _eventNonce,
        bool _success,
        bytes _cause
    );
}
