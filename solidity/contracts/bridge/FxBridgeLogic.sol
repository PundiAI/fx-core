// SPDX-License-Identifier: MIT

pragma experimental ABIEncoderV2;
pragma solidity ^0.8.0;

import {SafeMathUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/math/SafeMathUpgradeable.sol";
import {IERC20MetadataUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20MetadataUpgradeable.sol";
import {SafeERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/utils/SafeERC20Upgradeable.sol";
import {AddressUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/AddressUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";

import {IERC20ExtensionsUpgradeable} from "./IERC20ExtensionsUpgradeable.sol";

/* solhint-disable custom-errors */

contract FxBridgeLogic is
    ReentrancyGuardUpgradeable,
    OwnableUpgradeable,
    PausableUpgradeable
{
    using SafeMathUpgradeable for uint256;
    using SafeERC20Upgradeable for IERC20MetadataUpgradeable;
    using AddressUpgradeable for address;

    /* solhint-disable var-name-mixedcase */
    bytes32 public state_fxBridgeId;
    uint256 public state_powerThreshold;

    uint256 public state_lastEventNonce;
    bytes32 public state_lastOracleSetCheckpoint;
    uint256 public state_lastOracleSetNonce;
    mapping(address => uint256) public state_lastBatchNonces;

    address[] public bridgeTokens;
    mapping(address => TokenStatus) public tokenStatus;
    string public version;
    mapping(uint256 => bool) public state_lastRefundNonce;
    mapping(uint256 => bool) public state_lastBridgeCallNonce;
    /* solhint-enable var-name-mixedcase */

    struct TokenStatus {
        bool isOriginated;
        bool isActive;
        bool isExist;
        BridgeTokenType tokenType;
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
        BridgeTokenType tokenType;
    }

    struct BridgeCallData {
        address sender;
        address receiver;
        address to;
        uint256 value;
        bytes asset;
        bytes message;
        uint256 timeout;
        uint256 gasLimit;
    }

    enum BridgeTokenType {
        ERC20,
        ERC721,
        ERC404
    }
    /* =============== INIT =============== */

    function init(
        bytes32 _fxBridgeId,
        uint256 _powerThreshold,
        address[] memory _oracles,
        uint256[] memory _powers
    ) public initializer {
        __Pausable_init();
        __Ownable_init();
        __ReentrancyGuard_init();
        require(
            _oracles.length == _powers.length,
            "Malformed current oracle set"
        );

        uint256 cumulativePower = 0;
        for (uint256 i = 0; i < _powers.length; i++) {
            cumulativePower = cumulativePower + _powers[i];
            if (cumulativePower > _powerThreshold) {
                break;
            }
        }
        require(
            cumulativePower > _powerThreshold,
            "Submitted oracle set signatures do not have enough power."
        );

        bytes32 newCheckpoint = makeCheckpoint(
            _oracles,
            _powers,
            0,
            _fxBridgeId
        );

        state_fxBridgeId = _fxBridgeId;
        state_powerThreshold = _powerThreshold;
        state_lastOracleSetCheckpoint = newCheckpoint;
        state_lastOracleSetNonce = 0;
        state_lastEventNonce = 1;
        version = "1.0.0";

        emit OracleSetUpdatedEvent(
            state_lastOracleSetNonce,
            state_lastEventNonce,
            _oracles,
            _powers
        );
    }

    /* =============== MUTATIVE FUNCTIONS  =============== */

    function addBridgeToken(
        address _tokenAddr,
        bytes32 _channelIBC,
        bool _isOriginated,
        BridgeTokenType _tokenType
    ) public onlyOwner returns (bool) {
        require(_tokenAddr != address(0), "Invalid token address");
        require(
            tokenStatus[_tokenAddr].isExist == false,
            "Bridge token already exists"
        );
        _handlerAddBridgeToken(
            _tokenAddr,
            TokenStatus(_isOriginated, true, true, _tokenType)
        );
        emit AddBridgeTokenEvent(
            _tokenAddr,
            IERC20MetadataUpgradeable(_tokenAddr).name(),
            IERC20MetadataUpgradeable(_tokenAddr).symbol(),
            IERC20MetadataUpgradeable(_tokenAddr).decimals(),
            state_lastEventNonce,
            _channelIBC,
            _tokenType
        );
        return true;
    }

    function _handlerAddBridgeToken(
        address _tokenAddr,
        TokenStatus memory _tokenStatus
    ) internal {
        bridgeTokens.push(_tokenAddr);
        tokenStatus[_tokenAddr] = _tokenStatus;
        state_lastEventNonce = state_lastEventNonce.add(1);
    }

    function pauseBridgeToken(
        address _tokenAddr
    ) public onlyOwner returns (bool) {
        require(
            tokenStatus[_tokenAddr].isExist == true,
            "Bridge token doesn't exists"
        );
        require(
            tokenStatus[_tokenAddr].isActive == true,
            "Bridge token already paused"
        );
        tokenStatus[_tokenAddr].isActive = false;
        return true;
    }

    function activeBridgeToken(
        address _tokenAddr
    ) public onlyOwner returns (bool) {
        require(
            tokenStatus[_tokenAddr].isExist == true,
            "Bridge token doesn't exists"
        );
        require(
            tokenStatus[_tokenAddr].isActive == false,
            "Bridge token already actived"
        );
        tokenStatus[_tokenAddr].isActive = true;
        return true;
    }

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
    ) public whenNotPaused {
        require(
            _newOracleSetNonce > _currentOracleSetNonce,
            "New oracle set nonce must be greater than the current nonce"
        );

        require(
            _newOracles.length == _newPowers.length,
            "Malformed new oracle set"
        );

        require(
            _currentOracles.length == _currentPowers.length &&
                _currentOracles.length == _v.length &&
                _currentOracles.length == _r.length &&
                _currentOracles.length == _s.length,
            "Malformed current oracle set"
        );

        require(
            makeCheckpoint(
                _currentOracles,
                _currentPowers,
                _currentOracleSetNonce,
                state_fxBridgeId
            ) == state_lastOracleSetCheckpoint,
            "Supplied current oracles and powers do not match checkpoint."
        );

        bytes32 newCheckpoint = makeCheckpoint(
            _newOracles,
            _newPowers,
            _newOracleSetNonce,
            state_fxBridgeId
        );

        checkOracleSignatures(
            _currentOracles,
            _currentPowers,
            _v,
            _r,
            _s,
            newCheckpoint,
            state_powerThreshold
        );

        state_lastOracleSetCheckpoint = newCheckpoint;

        state_lastOracleSetNonce = _newOracleSetNonce;

        state_lastEventNonce = state_lastEventNonce.add(1);
        emit OracleSetUpdatedEvent(
            _newOracleSetNonce,
            state_lastEventNonce,
            _newOracles,
            _newPowers
        );
    }

    function sendToFx(
        address _tokenContract,
        bytes32 _destination,
        bytes32 _targetIBC,
        uint256 _amount
    ) public nonReentrant whenNotPaused {
        require(_amount > 0, "amount should be greater than zero");
        TokenStatus memory _tokenStatus = tokenStatus[_tokenContract];
        require(_tokenStatus.isExist, "Unsupported token address");
        require(_tokenStatus.isActive, "token was paused");
        require(
            _tokenStatus.tokenType == BridgeTokenType.ERC20,
            "Unsupported token type"
        );

        IERC20MetadataUpgradeable(_tokenContract).safeTransferFrom(
            msg.sender,
            address(this),
            _amount
        );
        if (_tokenStatus.isOriginated == true) {
            IERC20ExtensionsUpgradeable(_tokenContract).burn(_amount);
        }

        state_lastEventNonce = state_lastEventNonce.add(1);
        emit SendToFxEvent(
            _tokenContract,
            msg.sender,
            _destination,
            _targetIBC,
            _amount,
            state_lastEventNonce
        );
    }

    function bridgeCall(
        string memory _dstChainId,
        uint256 _gasLimit,
        address _receiver,
        address _to,
        bytes calldata _message,
        uint256 _value,
        bytes memory _asset
    ) external {
        // decode and transfer asset
        _transferAsset(_msgSender(), address(this), _asset, false);

        // last event nonce +1
        state_lastEventNonce = state_lastEventNonce.add(1);

        // bridge call event
        emit BridgeCallEvent(
            _msgSender(),
            _receiver,
            _to,
            state_lastEventNonce,
            _dstChainId,
            _gasLimit,
            _value,
            _message,
            _asset
        );
    }

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
    ) public nonReentrant whenNotPaused {
        {
            TokenStatus memory _tokenStatus = tokenStatus[_tokenContract];
            require(_tokenStatus.isExist, "Unsupported token address");
            require(_tokenStatus.isActive, "Token was paused");

            require(
                state_lastBatchNonces[_tokenContract] < _nonceArray[1],
                "New batch nonce must be greater than the current nonce."
            );

            require(
                block.number < _batchTimeout,
                "Batch timeout must be greater than the current block height."
            );

            require(
                _currentOracles.length == _currentPowers.length &&
                    _currentOracles.length == _v.length &&
                    _currentOracles.length == _r.length &&
                    _currentOracles.length == _s.length,
                "Malformed current oracle set."
            );

            require(
                makeCheckpoint(
                    _currentOracles,
                    _currentPowers,
                    _nonceArray[0],
                    state_fxBridgeId
                ) == state_lastOracleSetCheckpoint,
                "Supplied current oracles and powers do not match checkpoint."
            );

            require(
                _amounts.length == _destinations.length &&
                    _amounts.length == _fees.length,
                "Malformed batch of transactions."
            );

            checkOracleSignatures(
                _currentOracles,
                _currentPowers,
                _v,
                _r,
                _s,
                keccak256(
                    abi.encode(
                        state_fxBridgeId,
                        // bytes32 encoding of "transactionBatch"
                        0x7472616e73616374696f6e426174636800000000000000000000000000000000,
                        _amounts,
                        _destinations,
                        _fees,
                        _nonceArray[1],
                        _tokenContract,
                        _batchTimeout,
                        _feeReceive
                    )
                ),
                state_powerThreshold
            );

            state_lastBatchNonces[_tokenContract] = _nonceArray[1];

            {
                uint256 totalFee;
                for (uint256 i = 0; i < _amounts.length; i++) {
                    totalFee = totalFee.add(_fees[i]);
                    if (_tokenStatus.isOriginated == true) {
                        IERC20ExtensionsUpgradeable(_tokenContract).mint(
                            address(this),
                            _amounts[i]
                        );
                    }
                    IERC20MetadataUpgradeable(_tokenContract).safeTransfer(
                        _destinations[i],
                        _amounts[i]
                    );
                }

                if (_tokenStatus.isOriginated == true) {
                    IERC20ExtensionsUpgradeable(_tokenContract).mint(
                        address(this),
                        totalFee
                    );
                }
                IERC20MetadataUpgradeable(_tokenContract).safeTransfer(
                    _feeReceive,
                    totalFee
                );
            }
        }

        {
            state_lastEventNonce = state_lastEventNonce.add(1);
            emit TransactionBatchExecutedEvent(
                _nonceArray[1],
                _tokenContract,
                state_lastEventNonce
            );
        }
    }

    function refundBridgeToken(
        address[] memory _currentOracles,
        uint256[] memory _currentPowers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        uint256[2] memory _nonceArray,
        address _receiver,
        address[] memory _tokens,
        uint256[] memory _amounts,
        uint256 _timeout
    ) public nonReentrant whenNotPaused {
        {
            require(
                _tokens.length == _amounts.length,
                "Token not match amount"
            );

            for (uint256 i = 0; i < _tokens.length; i++) {
                TokenStatus memory _tokenStatus = tokenStatus[_tokens[i]];
                require(_tokenStatus.isExist, "Unsupported token address");
                require(_tokenStatus.isActive, "Token was paused");
                require(_amounts[i] > 0, "Amount must be great 0");
            }

            require(
                !state_lastRefundNonce[_nonceArray[1]],
                "New refund nonce must be not used."
            );

            require(
                block.number < _timeout,
                "refund timeout must be greater than the current block height."
            );

            require(
                _currentOracles.length == _currentPowers.length &&
                    _currentOracles.length == _v.length &&
                    _currentOracles.length == _r.length &&
                    _currentOracles.length == _s.length,
                "Malformed current oracle set."
            );

            require(
                makeCheckpoint(
                    _currentOracles,
                    _currentPowers,
                    _nonceArray[0],
                    state_fxBridgeId
                ) == state_lastOracleSetCheckpoint,
                "Supplied current oracles and powers do not match checkpoint."
            );

            checkOracleSignatures(
                _currentOracles,
                _currentPowers,
                _v,
                _r,
                _s,
                keccak256(
                    abi.encode(
                        state_fxBridgeId,
                        // bytes32 encoding of "refundToken"
                        0x726566756e64546f6b656e000000000000000000000000000000000000000000,
                        _receiver,
                        _tokens,
                        _amounts,
                        _nonceArray[1],
                        _timeout
                    )
                ),
                state_powerThreshold
            );
        }

        state_lastRefundNonce[_nonceArray[1]] = true;

        {
            for (uint256 i = 0; i < _tokens.length; i++) {
                TokenStatus memory _tokenStatus = tokenStatus[_tokens[i]];
                if (_tokenStatus.isOriginated == true) {
                    IERC20ExtensionsUpgradeable(_tokens[i]).mint(
                        address(this),
                        _amounts[i]
                    );
                }
                IERC20MetadataUpgradeable(_tokens[i]).safeTransfer(
                    _receiver,
                    _amounts[i]
                );
            }
        }

        {
            state_lastEventNonce = state_lastEventNonce.add(1);
            emit RefundTokenExecutedEvent(
                _receiver,
                _nonceArray[1],
                state_lastEventNonce
            );
        }
    }

    function submitBridgeCall(
        address[] memory _currentOracles,
        uint256[] memory _currentPowers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        uint256[2] memory _nonceArray,
        BridgeCallData memory _input
    ) public nonReentrant whenNotPaused {
        verifyBridgeCall(
            _currentOracles,
            _currentPowers,
            _v,
            _r,
            _s,
            _nonceArray,
            _input
        );

        state_lastBridgeCallNonce[_nonceArray[1]] = true;

        bool result = true;
        // solhint-disable-next-line no-empty-blocks
        try this.callAssetMessage(_input) {
            // event success
        } catch {
            // event failed
            result = false;
        }
        emit SubmitBridgeCallEvent(
            _input.sender,
            _input.receiver,
            _input.to,
            _nonceArray[1],
            false
        );
    }

    function bridgeCallSigHash(
        BridgeCallData memory input,
        uint256 nonce
    ) internal view returns (bytes32) {
        bytes memory data = abi.encode(
            state_fxBridgeId,
            // bytes32 encoding of "bridgeCall"
            0x62726964676543616c6c00000000000000000000000000000000000000000000,
            input.sender,
            input.receiver,
            input.to,
            input.value,
            input.asset,
            input.message,
            nonce,
            input.timeout
        );
        return keccak256(data);
    }

    function verifyBridgeCall(
        address[] memory _currentOracles,
        uint256[] memory _currentPowers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        uint256[2] memory _nonceArray,
        BridgeCallData memory _input
    ) internal view {
        require(
            !state_lastBridgeCallNonce[_nonceArray[1]],
            "New bridge call nonce must be not exist."
        );

        require(
            block.number < _input.timeout,
            "timeout must be greater than the current block height."
        );

        require(
            _currentOracles.length == _currentPowers.length &&
                _currentOracles.length == _v.length &&
                _currentOracles.length == _r.length &&
                _currentOracles.length == _s.length,
            "Malformed current oracle set."
        );

        require(
            makeCheckpoint(
                _currentOracles,
                _currentPowers,
                _nonceArray[0],
                state_fxBridgeId
            ) == state_lastOracleSetCheckpoint,
            "Supplied current oracles and powers do not match checkpoint."
        );

        bytes32 dataHash = bridgeCallSigHash(_input, _nonceArray[1]);

        checkOracleSignatures(
            _currentOracles,
            _currentPowers,
            _v,
            _r,
            _s,
            dataHash,
            state_powerThreshold
        );
    }

    function callAssetMessage(BridgeCallData memory _input) public onlySelf {
        if (_input.asset.length > 0) {
            _transferAsset(address(this), _input.receiver, _input.asset, true);
        }

        if (_input.message.length > 0) {
            bytes memory data = abi.encodeWithSignature(
                "onFxcoreMessage((address,address,address,uint256,bytes,bytes,uint256,uint256))",
                _input
            );
            // solhint-disable-next-line avoid-low-level-calls
            (bool success, ) = (_input.to).call{gas: _input.gasLimit}(data);
            require(success, "Call onFxcoreMessage failed");
        }
    }

    function transferOwner(
        address _token,
        address _newOwner
    ) public onlyOwner returns (bool) {
        IERC20ExtensionsUpgradeable(_token).transferOwnership(_newOwner);
        emit TransferOwnerEvent(_token, _newOwner);
        return true;
    }

    function updateTokenType(
        address[] memory _tokens,
        BridgeTokenType[] memory _tokenTypes
    ) public onlyOwner returns (bool) {
        for (uint256 i = 0; i < _tokens.length; i++) {
            address token = _tokens[i];
            BridgeTokenType tokenType = _tokenTypes[i];
            require(tokenStatus[token].isExist, "Unsupported token address");
            tokenStatus[token].tokenType = tokenType;
            state_lastEventNonce = state_lastEventNonce.add(1);
            emit AddBridgeTokenEvent(
                token,
                IERC20MetadataUpgradeable(token).name(),
                IERC20MetadataUpgradeable(token).symbol(),
                IERC20MetadataUpgradeable(token).decimals(),
                state_lastEventNonce,
                bytes32(0),
                tokenType
            );
        }
        return true;
    }

    /* =============== QUERY FUNCTIONS =============== */

    function lastBatchNonce(
        address _erc20Address
    ) public view returns (uint256) {
        return state_lastBatchNonces[_erc20Address];
    }

    function checkAssetStatus(address _tokenAddr) public view returns (bool) {
        return tokenStatus[_tokenAddr].isExist;
    }

    /* ============== HELP FUNCTIONS =============== */

    function verifySig(
        address _signer,
        bytes32 _theHash,
        uint8 _v,
        bytes32 _r,
        bytes32 _s
    ) private pure returns (bool) {
        bytes32 messageDigest = keccak256(
            abi.encodePacked("\x19Ethereum Signed Message:\n32", _theHash)
        );
        return _signer == ecrecover(messageDigest, _v, _r, _s);
    }

    function makeCheckpoint(
        address[] memory _oracles,
        uint256[] memory _powers,
        uint256 _oracleSetNonce,
        bytes32 _fxBridgeId
    ) public pure returns (bytes32) {
        // bytes32 encoding of the string "checkpoint"
        bytes32 methodName = 0x636865636b706f696e7400000000000000000000000000000000000000000000;
        return
            keccak256(
                abi.encode(
                    _fxBridgeId,
                    methodName,
                    _oracleSetNonce,
                    _oracles,
                    _powers
                )
            );
    }

    function checkOracleSignatures(
        address[] memory _currentOracles,
        uint256[] memory _currentPowers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        bytes32 _theHash,
        uint256 _powerThreshold
    ) public pure {
        uint256 cumulativePower = 0;

        for (uint256 i = 0; i < _currentOracles.length; i++) {
            if (_v[i] != 0) {
                require(
                    verifySig(
                        _currentOracles[i],
                        _theHash,
                        _v[i],
                        _r[i],
                        _s[i]
                    ),
                    "Oracle signature does not match."
                );
                cumulativePower = cumulativePower + _currentPowers[i];
                if (cumulativePower > _powerThreshold) {
                    break;
                }
            }
        }

        require(
            cumulativePower > _powerThreshold,
            "Submitted oracle set signatures do not have enough power."
        );
    }

    function pause() public onlyOwner whenNotPaused {
        _pause();
    }

    function unpause() public onlyOwner whenPaused {
        _unpause();
    }

    function getBridgeTokenList() public view returns (BridgeToken[] memory) {
        BridgeToken[] memory result = new BridgeToken[](bridgeTokens.length);
        for (uint256 i = 0; i < bridgeTokens.length; i++) {
            address _tokenAddr = address(bridgeTokens[i]);
            string memory _name = IERC20MetadataUpgradeable(_tokenAddr).name();
            string memory _symbol = IERC20MetadataUpgradeable(_tokenAddr)
                .symbol();
            uint8 _decimals = 0;
            if (tokenStatus[_tokenAddr].tokenType != BridgeTokenType.ERC721) {
                _decimals = IERC20MetadataUpgradeable(_tokenAddr).decimals();
            }
            BridgeToken memory bridgeToken = BridgeToken(
                _tokenAddr,
                _name,
                _symbol,
                _decimals,
                tokenStatus[_tokenAddr].tokenType
            );
            result[i] = bridgeToken;
        }
        return result;
    }

    function _transferAsset(
        address _from,
        address _receiver,
        bytes memory _asset,
        bool _mintToken
    ) internal {
        (address[] memory token, uint256[] memory amount) = _decodeAsset(
            _asset
        );
        for (uint256 i = 0; i < token.length; i++) {
            require(amount[i] > 0, "amount should be greater than zero");
            TokenStatus memory _tokenStatus = tokenStatus[token[i]];
            require(_tokenStatus.isExist, "Unsupported token address");
            require(_tokenStatus.isActive, "token was paused");

            // mint origin token
            if (_tokenStatus.isOriginated == true && _mintToken) {
                _mintAssetToken(token[i], _from, amount[i], _tokenStatus);
            }

            _transferAssetToken(
                token[i],
                _from,
                _receiver,
                amount[i],
                _tokenStatus
            );

            // burn origin token
            if (_tokenStatus.isOriginated == true && !_mintToken) {
                _burnAssetToken(token[i], amount[i], _tokenStatus);
            }
        }
    }

    function _decodeAsset(
        bytes memory _data
    ) internal pure returns (address[] memory, uint256[] memory) {
        (bytes memory tokenBytes, uint256[] memory amounts) = abi.decode(
            _data,
            (bytes, uint256[])
        );
        uint256 tokenCount = amounts.length;
        require(tokenBytes.length == tokenCount * 20, "Token not match amount");
        address[] memory tokens = new address[](tokenCount);
        for (uint256 i = 0; i < tokenCount; i++) {
            uint256 currentTokenStartingByte = 20 + i * 20;
            address currentToken;
            // solhint-disable-next-line no-inline-assembly
            assembly {
                currentToken := mload(add(tokenBytes, currentTokenStartingByte))
            }
            tokens[i] = currentToken;
        }
        return (tokens, amounts);
    }

    function _mintAssetToken(
        address _token,
        address _from,
        uint256 _amount,
        TokenStatus memory _tokenStatus
    ) internal {
        // ERC20
        if (_tokenStatus.tokenType == BridgeTokenType.ERC20) {
            IERC20ExtensionsUpgradeable(_token).mint(_from, _amount);
        }
        // todo ERC721 ERC404
    }

    function _burnAssetToken(
        address _token,
        uint256 _amount,
        TokenStatus memory _tokenStatus
    ) internal {
        // ERC20
        if (_tokenStatus.tokenType == BridgeTokenType.ERC20) {
            IERC20ExtensionsUpgradeable(_token).burn(_amount);
        }
        // todo ERC721 ERC404
    }

    function _transferAssetToken(
        address _token,
        address _from,
        address _to,
        uint256 _amount,
        TokenStatus memory _tokenStatus
    ) internal {
        if (_from == address(this)) {
            // ERC20
            if (_tokenStatus.tokenType == BridgeTokenType.ERC20) {
                IERC20MetadataUpgradeable(_token).safeTransfer(_to, _amount);
            }
            // todo ERC721 ERC404
        } else {
            // ERC20
            if (_tokenStatus.tokenType == BridgeTokenType.ERC20) {
                IERC20MetadataUpgradeable(_token).safeTransferFrom(
                    _from,
                    _to,
                    _amount
                );
            }
            // todo ERC721 ERC404
        }
    }

    modifier onlySelf() {
        require(
            address(this) == _msgSender(),
            "Selfable: caller is not the self"
        );
        _;
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
        bytes32 _channelIBC,
        BridgeTokenType _tokenType
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
        address indexed _receiver,
        address indexed _to,
        uint256 _eventNonce,
        string _dstChainId,
        uint256 _gasLimit,
        uint256 _value,
        bytes _message,
        bytes _asset
    );

    event RefundTokenExecutedEvent(
        address indexed _receiver,
        uint256 indexed _refundNonce,
        uint256 _eventNonce
    );

    event SubmitBridgeCallEvent(
        address indexed _sender,
        address indexed _receiver,
        address indexed _to,
        uint256 _eventNonce,
        bool _result
    );
}
