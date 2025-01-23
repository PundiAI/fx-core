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

import {IERC20ExtensionsUpgradeable} from "../interfaces/IERC20ExtensionsUpgradeable.sol";
import {IBridgeCallContext} from "../interfaces/IBridgeCallContext.sol";
import {IFxBridgeLogic} from "../interfaces/IFxBridgeLogic.sol";

/* solhint-disable custom-errors */

contract FxBridgeLogic is
    ReentrancyGuardUpgradeable,
    OwnableUpgradeable,
    PausableUpgradeable,
    IFxBridgeLogic
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
    mapping(uint256 => bool) public state_lastBridgeCallNonces;
    /* solhint-enable var-name-mixedcase */

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
            "Malformed current oracle set."
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

        state_fxBridgeId = _fxBridgeId;
        state_lastOracleSetCheckpoint = _oracleSetCheckpoint(
            0,
            _oracles,
            _powers
        );
        state_powerThreshold = _powerThreshold;
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
        bytes32 _memo,
        bool _isOriginated
    ) public onlyOwner returns (bool) {
        require(_tokenAddr != address(0), "Invalid token address.");
        require(
            tokenStatus[_tokenAddr].isExist == false,
            "Bridge token already exists."
        );
        bridgeTokens.push(_tokenAddr);
        tokenStatus[_tokenAddr] = TokenStatus(_isOriginated, true, true);
        state_lastEventNonce = state_lastEventNonce.add(1);
        emit AddBridgeTokenEvent(
            _tokenAddr,
            IERC20MetadataUpgradeable(_tokenAddr).name(),
            IERC20MetadataUpgradeable(_tokenAddr).symbol(),
            IERC20MetadataUpgradeable(_tokenAddr).decimals(),
            state_lastEventNonce,
            _memo
        );
        return true;
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
            "Bridge token doesn't exists."
        );
        require(
            tokenStatus[_tokenAddr].isActive == false,
            "Bridge token already activated."
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
            "New oracle set nonce must be greater than the current nonce."
        );

        require(
            _newOracles.length == _newPowers.length,
            "Malformed new oracle set."
        );

        OracleSignatures memory curOracleSigns = OracleSignatures(
            _currentOracles,
            _currentPowers,
            _r,
            _s,
            _v
        );

        _validateOracleSet(_currentOracleSetNonce, curOracleSigns);

        bytes32 newCheckpoint = _oracleSetCheckpoint(
            _newOracleSetNonce,
            _newOracles,
            _newPowers
        );

        _verifyOracleSig(curOracleSigns, newCheckpoint);

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
        require(_amount > 0, "Amount should be greater than zero.");
        TokenStatus memory _tokenStatus = tokenStatus[_tokenContract];
        require(_tokenStatus.isExist, "Unsupported token address.");
        require(_tokenStatus.isActive, "Token was paused.");

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
        string memory _dstChain,
        address _refund,
        address[] memory _tokens,
        uint256[] memory _amounts,
        address _to,
        bytes memory _data,
        uint256 _quoteId,
        uint256 _gasLimit,
        bytes memory _memo
    ) external payable nonReentrant whenNotPaused returns (uint256) {
        require(bytes(_dstChain).length == 0, "Invalid dstChain.");
        if (_tokens.length > 0) {
            require(_refund != address(0), "Refund address is empty.");
        }

        // transfer ERC20
        _transferERC20(_msgSender(), address(this), _tokens, _amounts);

        // last event nonce +1
        state_lastEventNonce = state_lastEventNonce.add(1);

        emit BridgeCallEvent(
            _msgSender(),
            _refund,
            _to,
            // solhint-disable-next-line avoid-tx-origin
            tx.origin,
            state_lastEventNonce,
            _dstChain,
            _tokens,
            _amounts,
            _data,
            _quoteId,
            _gasLimit,
            _memo
        );

        return state_lastEventNonce;
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
                _amounts.length == _destinations.length &&
                    _amounts.length == _fees.length,
                "Malformed batch of transactions."
            );

            OracleSignatures memory curOracleSigns = OracleSignatures(
                _currentOracles,
                _currentPowers,
                _r,
                _s,
                _v
            );

            _validateOracleSet(_nonceArray[0], curOracleSigns);

            _verifyOracleSig(
                curOracleSigns,
                submitBatchCheckpoint(
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

    function submitBridgeCall(
        OracleSignatures calldata _curOracleSigns,
        uint256[2] calldata _nonceArray,
        BridgeCallData calldata _input
    ) public nonReentrant whenNotPaused {
        _verifySubmitBridgeCall(_curOracleSigns, _nonceArray, _input);

        state_lastBridgeCallNonces[_nonceArray[1]] = true;

        bool success = false;
        bytes memory cause = new bytes(0);
        try this._onBridgeCall(_input) {
            success = true;
            // solhint-disable-next-line no-empty-blocks
        } catch Error(string memory reason) {
            // catch failing revert() and require()
            cause = bytes(reason);
        } catch (bytes memory reason) {
            // catch failing assert()
            cause = reason;
        }

        {
            state_lastEventNonce = state_lastEventNonce.add(1);
            emit SubmitBridgeCallEvent(
                // solhint-disable-next-line avoid-tx-origin
                tx.origin,
                _nonceArray[1],
                state_lastEventNonce,
                success,
                cause
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

    function getBridgeTokenList() public view returns (BridgeToken[] memory) {
        BridgeToken[] memory result = new BridgeToken[](bridgeTokens.length);
        for (uint256 i = 0; i < bridgeTokens.length; i++) {
            address _tokenAddr = address(bridgeTokens[i]);
            string memory _name = IERC20MetadataUpgradeable(_tokenAddr).name();
            string memory _symbol = IERC20MetadataUpgradeable(_tokenAddr)
                .symbol();
            uint8 _decimals = IERC20MetadataUpgradeable(_tokenAddr).decimals();
            BridgeToken memory bridgeToken = BridgeToken(
                _tokenAddr,
                _name,
                _symbol,
                _decimals
            );
            result[i] = bridgeToken;
        }
        return result;
    }

    function getTokenStatus(
        address _tokenAddr
    ) public view returns (TokenStatus memory) {
        return tokenStatus[_tokenAddr];
    }

    /* solhint-disable func-name-mixedcase */
    function convert_decimals(address) public pure returns (uint8) {
        revert("Not implemented");
    }

    /* ============== HELP FUNCTIONS =============== */

    function oracleSetCheckpoint(
        bytes32 _fxBridgeId,
        bytes32 _methodName,
        uint256 _oracleSetNonce,
        address[] memory _oracles,
        uint256[] memory _powers
    ) public pure returns (bytes32) {
        return
            keccak256(
                abi.encode(
                    _fxBridgeId,
                    _methodName,
                    _oracleSetNonce,
                    _oracles,
                    _powers
                )
            );
    }

    function submitBatchCheckpoint(
        bytes32 _fxBridgeId,
        bytes32 _methodName,
        uint256[] memory _amounts,
        address[] memory _destinations,
        uint256[] memory _fees,
        uint256 _batchNonce,
        address _tokenContract,
        uint256 _batchTimeout,
        address _feeReceive
    ) public pure returns (bytes32) {
        return
            keccak256(
                abi.encode(
                    _fxBridgeId,
                    _methodName,
                    _amounts,
                    _destinations,
                    _fees,
                    _batchNonce,
                    _tokenContract,
                    _batchTimeout,
                    _feeReceive
                )
            );
    }

    function bridgeCallCheckpoint(
        bytes32 _fxBridgeId,
        bytes32 _methodName,
        uint256 _nonce,
        BridgeCallData calldata _input
    ) public pure returns (bytes32) {
        return keccak256(abi.encode(_fxBridgeId, _methodName, _nonce, _input));
    }

    function _oracleSetCheckpoint(
        uint256 _oracleSetNonce,
        address[] memory _oracles,
        uint256[] memory _powers
    ) private view returns (bytes32) {
        return
            oracleSetCheckpoint(
                state_fxBridgeId,
                // bytes32 encoding of the string "checkpoint"
                0x636865636b706f696e7400000000000000000000000000000000000000000000,
                _oracleSetNonce,
                _oracles,
                _powers
            );
    }

    function _verifySubmitBridgeCall(
        OracleSignatures calldata _curOracleSigns,
        uint256[2] calldata _nonceArray,
        BridgeCallData calldata _input
    ) internal view {
        require(
            !state_lastBridgeCallNonces[_nonceArray[1]],
            "New bridge call nonce must be not exist."
        );

        require(
            block.number < _input.timeout,
            "Timeout must be greater than the current block height."
        );

        require(
            _input.tokens.length == _input.amounts.length,
            "Token not match amount."
        );

        _validateOracleSet(_nonceArray[0], _curOracleSigns);

        _verifyOracleSig(
            _curOracleSigns,
            bridgeCallCheckpoint(
                state_fxBridgeId,
                // bytes32 encoding of "bridgeCall"
                0x62726964676543616c6c00000000000000000000000000000000000000000000,
                _nonceArray[1],
                _input
            )
        );
    }

    function _onBridgeCall(BridgeCallData calldata _input) public onlySelf {
        if (_input.tokens.length > 0) {
            _transferERC20(
                address(this),
                _input.to,
                _input.tokens,
                _input.amounts
            );
        }

        if (_input.to.isContract()) {
            if (_input.eventNonce > 0) {
                IBridgeCallContext(_input.to).onRevert(
                    _input.eventNonce,
                    _input.data
                );
                return;
            }

            IBridgeCallContext(_input.to).onBridgeCall(
                _input.sender,
                _input.refund,
                _input.tokens,
                _input.amounts,
                _input.data,
                _input.memo
            );
        }
    }

    function _verifyOracleSig(
        OracleSignatures memory _curOracleSigns,
        bytes32 _theHash
    ) private view {
        uint256 cumulativePower = 0;

        for (uint256 i = 0; i < _curOracleSigns.oracles.length; i++) {
            if (_curOracleSigns.v[i] != 0) {
                require(
                    _verifySig(
                        _curOracleSigns.oracles[i],
                        _theHash,
                        _curOracleSigns.v[i],
                        _curOracleSigns.r[i],
                        _curOracleSigns.s[i]
                    ),
                    "Oracle signature does not match."
                );
                cumulativePower = cumulativePower + _curOracleSigns.powers[i];
                if (cumulativePower > state_powerThreshold) {
                    break;
                }
            }
        }

        require(
            cumulativePower > state_powerThreshold,
            "Submitted oracle set signatures do not have enough power."
        );
    }

    function _verifySig(
        address _signer,
        bytes32 _theHash,
        uint8 _v,
        bytes32 _r,
        bytes32 _s
    ) private pure returns (bool) {
        return
            _signer ==
            ecrecover(
                keccak256(
                    abi.encodePacked(
                        "\x19Ethereum Signed Message:\n32",
                        _theHash
                    )
                ),
                _v,
                _r,
                _s
            );
    }

    function _transferERC20(
        address _from,
        address _receiver,
        address[] memory _tokens,
        uint256[] memory _amounts
    ) internal {
        require(
            _tokens.length == _amounts.length,
            "Tokens and amounts not matched"
        );

        for (uint256 i = 0; i < _tokens.length; i++) {
            require(_tokens[i] != address(0), "Invalid token address");
            require(_amounts[i] > 0, "amount should be greater than zero");
            TokenStatus memory _tokenStatus = tokenStatus[_tokens[i]];
            require(_tokenStatus.isExist, "Unsupported token address");
            require(_tokenStatus.isActive, "token was paused");

            // mint origin token
            if (_tokenStatus.isOriginated == true && _from == address(this)) {
                IERC20ExtensionsUpgradeable(_tokens[i]).mint(
                    _from,
                    _amounts[i]
                );
            }

            if (_from == address(this)) {
                IERC20MetadataUpgradeable(_tokens[i]).safeTransfer(
                    _receiver,
                    _amounts[i]
                );
            } else {
                IERC20MetadataUpgradeable(_tokens[i]).safeTransferFrom(
                    _from,
                    _receiver,
                    _amounts[i]
                );
            }

            // burn origin token
            if (
                _tokenStatus.isOriginated == true && _receiver == address(this)
            ) {
                IERC20ExtensionsUpgradeable(_tokens[i]).burn(_amounts[i]);
            }
        }
    }

    function _validateOracleSet(
        uint256 _currentOracleSetNonce,
        OracleSignatures memory _curOracleSigns
    ) internal view {
        require(
            _curOracleSigns.oracles.length == _curOracleSigns.powers.length &&
                _curOracleSigns.oracles.length == _curOracleSigns.v.length &&
                _curOracleSigns.oracles.length == _curOracleSigns.r.length &&
                _curOracleSigns.oracles.length == _curOracleSigns.s.length,
            "Malformed current oracle set: array length mismatch."
        );

        require(
            _oracleSetCheckpoint(
                _currentOracleSetNonce,
                _curOracleSigns.oracles,
                _curOracleSigns.powers
            ) == state_lastOracleSetCheckpoint,
            "Supplied current oracles and powers do not match checkpoint."
        );
    }

    function pause() public onlyOwner whenNotPaused {
        _pause();
    }

    function unpause() public onlyOwner whenPaused {
        _unpause();
    }

    modifier onlySelf() {
        require(
            address(this) == _msgSender(),
            "Selfable: caller is not the self"
        );
        _;
    }
}
