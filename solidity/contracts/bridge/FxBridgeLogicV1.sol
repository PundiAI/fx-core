pragma experimental ABIEncoderV2;
pragma solidity ^0.6.6;

import "@openzeppelin/contracts-upgradeable/utils/math/SafeMathUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20MetadataUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/utils/SafeERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/AddressUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";

contract FxBridgeLogic is
    ReentrancyGuardUpgradeable,
    OwnableUpgradeable,
    PausableUpgradeable
{
    using SafeMathUpgradeable for uint256;
    using SafeERC20Upgradeable for IERC20MetadataUpgradeable;

    bytes32 public state_fxBridgeId;
    uint256 public state_powerThreshold;

    address public state_fxOriginatedToken;
    uint256 public state_lastEventNonce;
    bytes32 public state_lastValsetCheckpoint;
    uint256 public state_lastValsetNonce;
    mapping(address => uint256) public state_lastBatchNonces;
    address[] public bridgeTokens;

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
        address[] memory _validators,
        uint256[] memory _powers
    ) public initializer {
        __Pausable_init();
        __Ownable_init();
        __ReentrancyGuard_init();
        require(
            _validators.length == _powers.length,
            "Malformed current validator set"
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
            "Submitted validator set signatures do not have enough power."
        );

        bytes32 newCheckpoint = makeCheckpoint(
            _validators,
            _powers,
            0,
            _fxBridgeId
        );

        state_fxBridgeId = _fxBridgeId;
        state_powerThreshold = _powerThreshold;
        state_lastValsetCheckpoint = newCheckpoint;
        state_lastValsetNonce = 0;
        state_lastEventNonce = 1;

        emit ValsetUpdatedEvent(
            state_lastValsetNonce,
            state_lastEventNonce,
            _validators,
            _powers
        );
    }

    /* =============== MUTATIVE FUNCTIONS  =============== */

    function setFxOriginatedToken(
        address _tokenAddr
    ) public onlyOwner returns (bool) {
        require(_tokenAddr != state_fxOriginatedToken, 'Invalid bridge token');
        state_fxOriginatedToken = _tokenAddr;
        state_lastEventNonce = state_lastEventNonce.add(1);
        emit FxOriginatedTokenEvent(
            _tokenAddr,
            IERC20MetadataUpgradeable(_tokenAddr).name(),
            IERC20MetadataUpgradeable(_tokenAddr).symbol(),
            IERC20MetadataUpgradeable(_tokenAddr).decimals(),
            state_lastEventNonce
        );
        return true;
    }

    function addBridgeToken(
        address _tokenAddr
    ) public onlyOwner returns (bool) {
        require(_tokenAddr != address(0), "Invalid address");
        require(_tokenAddr != state_fxOriginatedToken, 'Invalid bridge token');
        require(!_isContainToken(bridgeTokens, _tokenAddr), "Token already exists");
        bridgeTokens.push(_tokenAddr);
        return true;
    }

    function delBridgeToken(address _tokenAddr) public onlyOwner returns (bool) {
        require(_isContainToken(bridgeTokens, _tokenAddr), "Token not exists");
        for (uint i = 0; i < bridgeTokens.length; i++) {
            if (_tokenAddr == bridgeTokens[i]) {
                bridgeTokens[i] = bridgeTokens[bridgeTokens.length - 1];
                bridgeTokens.pop();
                return true;
            }
        }
        return false;
    }

    function updateValset(
        address[] memory _newValidators,
        uint256[] memory _newPowers,
        uint256 _newValsetNonce,
        address[] memory _currentValidators,
        uint256[] memory _currentPowers,
        uint256 _currentValsetNonce,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s
    ) public {

        require(
            _newValsetNonce > _currentValsetNonce,
            "New valset nonce must be greater than the current nonce"
        );

        require(
            _newValidators.length == _newPowers.length,
            "Malformed new validator set"
        );

        require(
            _currentValidators.length == _currentPowers.length &&
            _currentValidators.length == _v.length &&
            _currentValidators.length == _r.length &&
            _currentValidators.length == _s.length,
            "Malformed current validator set"
        );

        require(
            makeCheckpoint(
                _currentValidators,
                _currentPowers,
                _currentValsetNonce,
                state_fxBridgeId
            ) == state_lastValsetCheckpoint,
            "Supplied current validators and powers do not match checkpoint."
        );

        bytes32 newCheckpoint = makeCheckpoint(
            _newValidators,
            _newPowers,
            _newValsetNonce,
            state_fxBridgeId
        );

        checkValidatorSignatures(
            _currentValidators,
            _currentPowers,
            _v,
            _r,
            _s,
            newCheckpoint,
            state_powerThreshold
        );

        state_lastValsetCheckpoint = newCheckpoint;

        state_lastValsetNonce = _newValsetNonce;

        state_lastEventNonce = state_lastEventNonce.add(1);
        emit ValsetUpdatedEvent(
            _newValsetNonce,
            state_lastEventNonce,
            _newValidators,
            _newPowers
        );
    }

    function sendToFx(
        address _tokenContract,
        bytes32 _destination,
        bytes32 _targetIBC,
        uint256 _amount
    ) public nonReentrant whenNotPaused {
        require(checkAssetStatus(_tokenContract), "Unsupported token address");

        IERC20MetadataUpgradeable(_tokenContract).transferFrom(
            msg.sender,
            address(this),
            _amount
        );
        if (_tokenContract == state_fxOriginatedToken) {
            IERC20MetadataUpgradeable(_tokenContract).burn(_amount);
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
        address[] memory _currentValidators,
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
    ) public nonReentrant {
        {
            require(checkAssetStatus(_tokenContract), "Unsupported token address.");

            require(
                state_lastBatchNonces[_tokenContract] < _nonceArray[1],
                "New batch nonce must be greater than the current nonce."
            );

            require(
                block.number < _batchTimeout,
                "Batch timeout must be greater than the current block height."
            );

            require(
                _currentValidators.length == _currentPowers.length &&
                _currentValidators.length == _v.length &&
                _currentValidators.length == _r.length &&
                _currentValidators.length == _s.length,
                "Malformed current validator set."
            );

            require(
                makeCheckpoint(
                    _currentValidators,
                    _currentPowers,
                    _nonceArray[0],
                    state_fxBridgeId
                ) == state_lastValsetCheckpoint,
                "Supplied current validators and powers do not match checkpoint."
            );

            require(
                _amounts.length == _destinations.length &&
                    _amounts.length == _fees.length,
                "Malformed batch of transactions."
            );

            checkValidatorSignatures(
                _currentValidators,
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
                    if (_tokenContract == state_fxOriginatedToken) {
                        IERC20MetadataUpgradeable(state_fxOriginatedToken).mint(
                            _destinations[i],
                            _amounts[i]
                        );
                    } else {
                        IERC20MetadataUpgradeable(_tokenContract).safeTransfer(
                            _destinations[i],
                            _amounts[i]
                        );
                    }
                    totalFee = totalFee.add(_fees[i]);
                }

                if (_tokenContract == state_fxOriginatedToken) {
                    IERC20MetadataUpgradeable(state_fxOriginatedToken).mint(
                        _feeReceive,
                        totalFee
                    );
                } else {
                    IERC20MetadataUpgradeable(_tokenContract).safeTransfer(
                        _feeReceive,
                        totalFee
                    );
                }
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
        IERC20MetadataUpgradeable(_token).transferOwnership(_newOwner);
        emit TransferOwnerEvent(_token, _newOwner);
        return true;
    }

    /* =============== QUERY FUNCTIONS =============== */

    function lastBatchNonce(
        address _erc20Address
    ) public view returns (uint256) {
        return state_lastBatchNonces[_erc20Address];
    }

    function getBridgeTokenList() public view returns (BridgeToken[] memory) {
        BridgeToken[] memory result = new BridgeToken[](bridgeTokens.length);
        for (uint256 i = 0; i < bridgeTokens.length; i++) {
            address _tokenAddr = address(bridgeTokens[i]);
            BridgeToken memory bridgeToken = BridgeToken(
                _tokenAddr,
                IERC20MetadataUpgradeable(_tokenAddr).name(),
                IERC20MetadataUpgradeable(_tokenAddr).symbol(),
                IERC20MetadataUpgradeable(_tokenAddr).decimals());
            result[i] = bridgeToken;
        }
        return result;
    }

    function checkAssetStatus(address _tokenAddr) public view returns (bool){
        if (state_fxOriginatedToken == _tokenAddr) {
            return true;
        }
        return _isContainToken(bridgeTokens, _tokenAddr);
    }

    /* ============== HELP FUNCTIONS =============== */

    function _isContainToken(address[] memory list, address _tokenAddr) private pure returns (bool) {
        for (uint i = 0; i < list.length; i++) {
            if (list[i] == _tokenAddr) {
                return true;
            }
        }
        return false;
    }

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
        address[] memory _validators,
        uint256[] memory _powers,
        uint256 _valsetNonce,
        bytes32 _fxBridgeId
    ) public pure returns (bytes32) {
        // bytes32 encoding of the string "checkpoint"
        bytes32 methodName = 0x636865636b706f696e7400000000000000000000000000000000000000000000;
        return
            keccak256(
                abi.encode(
                    _fxBridgeId,
                    methodName,
                    _valsetNonce,
                    _validators,
                    _powers
                )
            );
    }

    function checkValidatorSignatures(
        address[] memory _currentValidators,
        uint256[] memory _currentPowers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        bytes32 _theHash,
        uint256 _powerThreshold
    ) public pure {
        uint256 cumulativePower = 0;

        for (uint256 i = 0; i < _currentValidators.length; i++) {
            if (_v[i] != 0) {
                require(
                    verifySig(
                        _currentValidators[i],
                        _theHash,
                        _v[i],
                        _r[i],
                        _s[i]
                    ),
                    "Validator signature does not match."
                );
                cumulativePower = cumulativePower + _currentPowers[i];
                if (cumulativePower > _powerThreshold) {
                    break;
                }
            }
        }

        require(
            cumulativePower > _powerThreshold,
            "Submitted validator set signatures do not have enough power."
        );
    }

    function pause() public onlyOwner whenNotPaused {
        _pause();
    }

    function unpause() public onlyOwner whenPaused {
        _unpause();
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
    event FxOriginatedTokenEvent(
        address indexed _tokenContract,
        string _name,
        string _symbol,
        uint8 _decimals,
        uint256 _eventNonce
    );
    event ValsetUpdatedEvent(
        uint256 indexed _newValsetNonce,
        uint256 _eventNonce,
        address[] _validators,
        uint256[] _powers
    );
    event TransferOwnerEvent(address _token, address _newOwner);
}