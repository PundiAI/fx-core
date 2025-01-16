// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {IBridgeFeeQuote, IBridgeFeeOracle} from "../interfaces/IBridgeFee.sol";
import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import {IERC20MetadataUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20MetadataUpgradeable.sol";
import {SafeERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/utils/SafeERC20Upgradeable.sol";
import {StringsUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/StringsUpgradeable.sol";
import {ECDSAUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/ECDSAUpgradeable.sol";
import {EnumerableSetUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/structs/EnumerableSetUpgradeable.sol";

error ChainNameInvalid();
error TokenNameInvalid();
error OracleInvalid();
error QuoteCapInvalid();
error ChainNameAlreadyExists();
error TokenNameAlreadyExists();

contract BridgeFeeQuote is
    IBridgeFeeQuote,
    Initializable,
    UUPSUpgradeable,
    AccessControlUpgradeable,
    ReentrancyGuardUpgradeable
{
    using EnumerableSetUpgradeable for EnumerableSetUpgradeable.Bytes32Set;

    using SafeERC20Upgradeable for IERC20MetadataUpgradeable;
    using ECDSAUpgradeable for bytes32;
    using StringsUpgradeable for string;

    bytes32 public constant OWNER_ROLE = keccak256("OWNER_ROLE");
    bytes32 public constant UPGRADE_ROLE = keccak256("UPGRADE_ROLE");

    uint8 public maxQuoteCap;
    uint256 public quoteNonce;
    address public oracleContract;

    EnumerableSetUpgradeable.Bytes32Set private chainNames;
    mapping(bytes32 => EnumerableSetUpgradeable.Bytes32Set) private tokens;

    // quotes is a mapping of quote id to quote
    mapping(uint256 => Quote) private quotes;
    // quoteIndexes is a mapping of quote index key to quote index
    mapping(bytes32 => QuoteIndex) private quoteIndexes;

    function initialize(
        address _oracle,
        uint8 _maxQuoteCap
    ) public initializer {
        oracleContract = _oracle;
        maxQuoteCap = _maxQuoteCap;
        quoteNonce = 1;

        __AccessControl_init();
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();

        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);
        _grantRole(OWNER_ROLE, msg.sender);
    }

    /**
     * @notice Quote the bridge fee for a given chainName, token and oracle.
     * @param _inputs QuoteInput[] The quote inputs.
     * @return uint256[] The quote ids.
     */
    function quote(
        QuoteInput[] calldata _inputs
    ) external nonReentrant returns (uint256[] memory) {
        uint256[] memory ids = new uint256[](_inputs.length);
        for (uint256 i = 0; i < _inputs.length; i++) {
            verifyInput(_inputs[i]);

            if (
                !IBridgeFeeOracle(oracleContract).isOnline(
                    _inputs[i].chainName,
                    msg.sender
                )
            ) {
                revert OracleInvalid();
            }

            bytes32 key = quoteIndexKey(
                _inputs[i].chainName,
                _inputs[i].tokenName,
                msg.sender,
                _inputs[i].cap
            );
            uint256 id = quoteIndexes[key].id;
            if (id != 0) {
                delete quotes[id];
                quoteIndexes[key].id = quoteNonce;
            } else {
                quoteIndexes[key] = QuoteIndex(
                    quoteNonce,
                    _inputs[i].chainName,
                    _inputs[i].tokenName,
                    msg.sender,
                    _inputs[i].cap
                );
            }
            quotes[quoteNonce] = Quote(
                _inputs[i].expiry,
                _inputs[i].gasLimit,
                _inputs[i].amount,
                key
            );
            ids[i] = quoteNonce;

            emit NewQuote(
                quoteNonce,
                _inputs[i].chainName,
                _inputs[i].tokenName,
                msg.sender,
                _inputs[i].amount,
                _inputs[i].gasLimit,
                _inputs[i].expiry,
                _inputs[i].cap
            );
            quoteNonce++;
        }
        return ids;
    }

    /**
     * @notice Register a new chain.
     * @param _chainName bytes32 The chain name.
     * @param _tokens bytes32[] The tokens.
     * @return bool Whether the chain is registered.
     */
    function registerChain(
        bytes32 _chainName,
        bytes32[] calldata _tokens
    ) external onlyRole(OWNER_ROLE) returns (bool) {
        if (chainNames.contains(_chainName)) {
            revert ChainNameAlreadyExists();
        }
        chainNames.add(_chainName);
        for (uint256 i = 0; i < _tokens.length; i++) {
            if (tokens[_chainName].contains(_tokens[i])) {
                revert TokenNameAlreadyExists();
            }
            tokens[_chainName].add(_tokens[i]);
        }
        return true;
    }

    /**
     * @notice Add new tokens to a chain.
     * @param _chainName bytes32 The chain name.
     * @param _tokens bytes32[] The tokens.
     * @return bool Whether the token is added.
     */
    function addToken(
        bytes32 _chainName,
        bytes32[] calldata _tokens
    ) external onlyRole(OWNER_ROLE) returns (bool) {
        if (!chainNames.contains(_chainName)) {
            revert ChainNameInvalid();
        }
        for (uint256 i = 0; i < _tokens.length; i++) {
            if (tokens[_chainName].contains(_tokens[i])) {
                revert TokenNameAlreadyExists();
            }
            tokens[_chainName].add(_tokens[i]);
        }
        return true;
    }

    /**
     * @notice Remove tokens from a chain.
     * @param _chainName bytes32 The chain name.
     * @param _tokens bytes32[] The tokens.
     * @return bool Whether the token is removed.
     */
    function removeToken(
        bytes32 _chainName,
        bytes32[] calldata _tokens
    ) external onlyRole(OWNER_ROLE) returns (bool) {
        for (uint256 i = 0; i < _tokens.length; i++) {
            if (!tokens[_chainName].contains(_tokens[i])) {
                revert TokenNameInvalid();
            }
            tokens[_chainName].remove(_tokens[i]);
        }
        return true;
    }

    function getQuoteById(
        uint256 _id
    ) external view returns (QuoteInfo memory) {
        Quote memory quoteInfo = quotes[_id];
        QuoteIndex memory index = quoteIndexes[quoteInfo.index];
        return
            QuoteInfo(
                _id,
                index.chainName,
                index.tokenName,
                index.oracle,
                quoteInfo.amount,
                quoteInfo.gasLimit,
                quoteInfo.expiry
            );
    }

    /**
     * @notice Get quotes by token.
     * @param _chainName bytes32 The chain name.
     * @param _token bytes32 The token.
     * @return QuoteInfo[] The quotes.
     */
    function getQuotesByToken(
        bytes32 _chainName,
        bytes32 _token
    )
        external
        view
        activeToken(_chainName, _token)
        returns (QuoteInfo[] memory)
    {
        address[] memory oracles = IBridgeFeeOracle(oracleContract)
            .getOracleList(_chainName);

        uint256 length = oracles.length * maxQuoteCap;

        QuoteInfo[] memory quotesList = new QuoteInfo[](length);

        uint256 index = 0;
        for (uint256 i = 0; i < oracles.length; i++) {
            for (uint8 j = 0; j < maxQuoteCap; j++) {
                quotesList[index] = getQuoteByIndex(
                    _chainName,
                    _token,
                    oracles[i],
                    j
                );
                index++;
            }
        }
        return quotesList;
    }

    /**
     * @notice Get default oracle quotes by token.
     * @param _chainName bytes32 The chain name.
     * @param _token bytes32 The token.
     * @return QuoteInfo[] The quotes.
     */
    function getDefaultOracleQuote(
        bytes32 _chainName,
        bytes32 _token
    )
        external
        view
        activeToken(_chainName, _token)
        returns (QuoteInfo[] memory)
    {
        address oracle = IBridgeFeeOracle(oracleContract).defaultOracle();

        QuoteInfo[] memory quotesList = new QuoteInfo[](maxQuoteCap);
        for (uint8 i = 0; i < maxQuoteCap; i++) {
            quotesList[i] = getQuoteByIndex(_chainName, _token, oracle, i);
        }

        return quotesList;
    }

    /**
     * @notice Get quote by index.
     * @param _chainName bytes32 The chain name.
     * @param _token bytes32 The token.
     * @param _oracle address The oracle.
     * @param _cap uint8 The cap.
     * @return QuoteInfo The quote.
     */
    function getQuoteByIndex(
        bytes32 _chainName,
        bytes32 _token,
        address _oracle,
        uint8 _cap
    ) public view activeToken(_chainName, _token) returns (QuoteInfo memory) {
        bytes32 key = quoteIndexKey(_chainName, _token, _oracle, _cap);
        QuoteIndex memory index = quoteIndexes[key];
        Quote memory quoteInfo = quotes[index.id];
        return
            QuoteInfo(
                index.id,
                index.chainName,
                index.tokenName,
                index.oracle,
                quoteInfo.amount,
                quoteInfo.gasLimit,
                quoteInfo.expiry
            );
    }

    function getChainNames() external view returns (bytes32[] memory) {
        return chainNames.values();
    }

    function getTokens(
        bytes32 _chainName
    ) external view returns (bytes32[] memory) {
        return tokens[_chainName].values();
    }

    function verifyInput(
        QuoteInput calldata _input
    )
        internal
        view
        activeToken(_input.chainName, _input.tokenName)
        returns (bool)
    {
        if (_input.cap >= maxQuoteCap) {
            revert QuoteCapInvalid();
        }
        return true;
    }

    modifier activeToken(bytes32 _chainName, bytes32 _token) {
        if (!chainNames.contains(_chainName)) {
            revert ChainNameInvalid();
        }
        if (!tokens[_chainName].contains(_token)) {
            revert TokenNameInvalid();
        }
        _;
    }

    function quoteIndexKey(
        bytes32 _chainName,
        bytes32 _token,
        address _oracle,
        uint8 _cap
    ) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(_chainName, _token, _oracle, _cap));
    }

    receive() external payable {}

    function _authorizeUpgrade(
        address
    ) internal override onlyRole(UPGRADE_ROLE) {} // solhint-disable-line no-empty-blocks
}
