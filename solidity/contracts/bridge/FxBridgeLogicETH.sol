// SPDX-License-Identifier: Apache-2.0

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

contract FxBridgeLogicETH is
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

    address public state_fxOriginatedToken;
    uint256 public state_lastEventNonce;
    bytes32 public state_lastOracleSetCheckpoint;
    uint256 public state_lastOracleSetNonce;
    mapping(address => uint256) public state_lastBatchNonces;
    /* solhint-enable var-name-mixedcase */

    address[] public bridgeTokens;
    /**
     * @dev Update params support fxCore v4
     *  Compatible with cross chain interfaces
     */

    mapping(address => TokenStatus) public tokenStatus;
    string public version;

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

    /* setFxOriginatedToken Deprecated: after fxcore upgrade v3 */
    function setFxOriginatedToken(
        address _tokenAddr
    ) public onlyOwner returns (bool) {
        require(_tokenAddr != state_fxOriginatedToken, "Invalid bridge token");
        state_fxOriginatedToken = _tokenAddr;
        state_lastEventNonce = state_lastEventNonce.add(1);
        emit FxOriginatedTokenEvent(
            _tokenAddr,
            IERC20MetadataUpgradeable(_tokenAddr).name(),
            IERC20MetadataUpgradeable(_tokenAddr).symbol(),
            IERC20MetadataUpgradeable(_tokenAddr).decimals(),
            state_lastEventNonce
        );
        emit AddBridgeTokenEvent(
            _tokenAddr,
            IERC20MetadataUpgradeable(_tokenAddr).name(),
            IERC20MetadataUpgradeable(_tokenAddr).symbol(),
            IERC20MetadataUpgradeable(_tokenAddr).decimals(),
            state_lastEventNonce,
            bytes32(0)
        );
        return true;
    }

    function addBridgeToken(
        address _tokenAddr,
        bytes32 _channelIBC,
        bool _isOriginated
    ) public onlyOwner returns (bool) {
        require(_tokenAddr != address(0), "Invalid token address");
        require(
            tokenStatus[_tokenAddr].isExist == false,
            "Bridge token already exists"
        );
        _handlerAddBridgeToken(
            _tokenAddr,
            TokenStatus(_isOriginated, true, true)
        );
        emit AddBridgeTokenEvent(
            _tokenAddr,
            IERC20MetadataUpgradeable(_tokenAddr).name(),
            IERC20MetadataUpgradeable(_tokenAddr).symbol(),
            IERC20MetadataUpgradeable(_tokenAddr).decimals(),
            state_lastEventNonce,
            _channelIBC
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

    function transferOwner(
        address _token,
        address _newOwner
    ) public onlyOwner returns (bool) {
        IERC20ExtensionsUpgradeable(_token).transferOwnership(_newOwner);
        emit TransferOwnerEvent(_token, _newOwner);
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
            BridgeToken memory bridgeToken = BridgeToken(
                _tokenAddr,
                IERC20MetadataUpgradeable(_tokenAddr).name(),
                IERC20MetadataUpgradeable(_tokenAddr).symbol(),
                IERC20MetadataUpgradeable(_tokenAddr).decimals()
            );
            result[i] = bridgeToken;
        }
        return result;
    }

    /* ============== UPGRADE FUNCTIONS =============== */

    function migrate() public onlyOwner {
        version = "1.1.0";

        for (uint256 i = 0; i < bridgeTokens.length; i++) {
            tokenStatus[address(bridgeTokens[i])] = TokenStatus(
                false,
                true,
                true
            );
        }
        bridgeTokens.push(state_fxOriginatedToken);
        tokenStatus[state_fxOriginatedToken] = TokenStatus(true, true, true);
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
    /* FxOriginatedTokenEvent Deprecated: after fxcore upgrade v3 */
    event FxOriginatedTokenEvent(
        address indexed _tokenContract,
        string _name,
        string _symbol,
        uint8 _decimals,
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
}
