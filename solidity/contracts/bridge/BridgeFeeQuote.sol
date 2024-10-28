// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {IBridgeFeeQuote, IBridgeFeeOracle} from "./IBridgeFee.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import {IERC20MetadataUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20MetadataUpgradeable.sol";
import {SafeERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/utils/SafeERC20Upgradeable.sol";
import {ECDSAUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/ECDSAUpgradeable.sol";
import {StringsUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/StringsUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

error ChainNameInvalid();
error TokenNameInvalid();
error OracleInvalid();
error QuoteIndexInvalid();
error QuoteExpired();
error VerifySignatureFailed(address, address);
error QuoteNotFound();
error ChainNameAlreadyExists();
error TokenNameAlreadyExists();

contract BridgeFeeQuote is
    IBridgeFeeQuote,
    Initializable,
    UUPSUpgradeable,
    OwnableUpgradeable,
    ReentrancyGuardUpgradeable
{
    using SafeERC20Upgradeable for IERC20MetadataUpgradeable;
    using ECDSAUpgradeable for bytes32;
    using StringsUpgradeable for string;

    struct Quote {
        uint256 id;
        uint256 fee;
        uint256 gasLimit;
        uint256 expiry;
    }

    address public oracleContract;
    // maximum number of quotes per oracle
    uint256 public maxQuoteIndex;

    uint256 public quoteNonce;

    string[] public chainNames;

    // chainName -> Assert
    mapping(string => Asset) public assets;

    // Only one quote is allowed per chainName + tokenName + oracle + index
    mapping(bytes => Quote) internal quotes; // key: chainName + tokenName + oracle + index
    // id -> chainName + tokenName + oracle + index
    mapping(uint256 => bytes) internal quoteIds;

    function initialize(
        address _oracleContract,
        uint256 _maxQuoteIndex
    ) public initializer {
        oracleContract = _oracleContract;
        maxQuoteIndex = _maxQuoteIndex;

        __Ownable_init();
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
    }

    event NewQuote(
        uint256 indexed id,
        address indexed oracle,
        string indexed chainName,
        string tokenName,
        uint256 fee,
        uint256 gasLimit,
        uint256 expiry
    );

    /**
     * @notice Quote the bridge fee for a given chainName, token and oracle.
     * @param _inputs QuoteInput[] The quote inputs.
     * @return bool Whether the quote is successful.
     */
    function quote(
        QuoteInput[] memory _inputs
    ) external nonReentrant returns (bool) {
        for (uint256 i = 0; i < _inputs.length; i++) {
            verifyInput(_inputs[i]);

            verifySignature(_inputs[i]);

            bytes memory asset = packAsset(
                _inputs[i].chainName,
                _inputs[i].tokenName,
                _inputs[i].oracle,
                _inputs[i].quoteIndex
            );
            if (quotes[asset].id > 0) {
                delete quoteIds[quotes[asset].id];
            }
            quotes[asset] = Quote({
                id: quoteNonce,
                fee: _inputs[i].fee,
                gasLimit: _inputs[i].gasLimit,
                expiry: _inputs[i].expiry
            });
            quoteIds[quoteNonce] = asset;
            emit NewQuote(
                quoteNonce,
                _inputs[i].oracle,
                _inputs[i].chainName,
                _inputs[i].tokenName,
                _inputs[i].fee,
                _inputs[i].gasLimit,
                _inputs[i].expiry
            );
            quoteNonce = quoteNonce + 1;
        }
        return true;
    }

    /**
     * @notice Get the quote list for a given chainName.
     * @param _chainName The name of the chain.
     * @return QuoteInfo[] The quote list.
     */
    function getQuoteList(
        string memory _chainName
    ) external view returns (QuoteInfo[] memory) {
        if (!assets[_chainName].isActive) {
            revert ChainNameInvalid();
        }

        QuoteInfo[] memory quotesList = new QuoteInfo[](
            currentActiveQuoteNum(_chainName)
        );
        uint256 currentIndex = 0;
        address[] memory oracles = IBridgeFeeOracle(oracleContract)
            .getOracleList();

        for (uint256 i = 0; i < oracles.length; i++) {
            for (uint256 j = 0; j < assets[_chainName].tokenNames.length; j++) {
                for (uint256 k = 0; k < maxQuoteIndex; k++) {
                    bytes memory asset = packAsset(
                        _chainName,
                        assets[_chainName].tokenNames[j],
                        oracles[i],
                        k
                    );
                    if (quotes[asset].expiry >= block.timestamp) {
                        quotesList[currentIndex] = QuoteInfo({
                            id: quotes[asset].id,
                            chainName: _chainName,
                            tokenName: assets[_chainName].tokenNames[j],
                            oracle: oracles[i],
                            fee: quotes[asset].fee,
                            gasLimit: quotes[asset].gasLimit,
                            expiry: quotes[asset].expiry
                        });
                        currentIndex += 1;
                    }
                }
            }
        }
        return quotesList;
    }

    /**
     * @notice Get the quote by id.
     * @param _id The id of the quote.
     * @return q QuoteInfo The quote.
     */
    function getQuoteById(
        uint256 _id
    ) external view returns (QuoteInfo memory q) {
        bytes memory asset = quoteIds[_id];
        if (asset.length == 0) {
            revert QuoteNotFound();
        }
        (
            string memory chainName,
            string memory tokenName,
            address oracle,

        ) = unpackAsset(asset);
        q = QuoteInfo({
            id: _id,
            chainName: chainName,
            tokenName: tokenName,
            oracle: oracle,
            fee: quotes[asset].fee,
            gasLimit: quotes[asset].gasLimit,
            expiry: quotes[asset].expiry
        });
    }

    /**
     * @notice Get the quotes by token.
     * @param _chainName The name of the chain.
     * @param _tokenName The address of the token.
     * @return QuoteInfo[] The quote list.
     */
    function getQuoteByToken(
        string memory _chainName,
        string memory _tokenName
    ) external view returns (QuoteInfo[] memory) {
        if (!assets[_chainName].isActive) {
            revert ChainNameInvalid();
        }
        address oracle = IBridgeFeeOracle(oracleContract).defaultOracle();

        QuoteInfo[] memory quotesList = new QuoteInfo[](maxQuoteIndex);

        for (uint256 i = 0; i < maxQuoteIndex; i++) {
            bytes memory asset = packAsset(_chainName, _tokenName, oracle, i);
            quotesList[i] = QuoteInfo({
                id: quotes[asset].id,
                chainName: _chainName,
                tokenName: _tokenName,
                oracle: oracle,
                fee: quotes[asset].fee,
                gasLimit: quotes[asset].gasLimit,
                expiry: quotes[asset].expiry
            });
        }
        return quotesList;
    }

    function currentActiveQuoteNum(
        string memory _chainName
    ) internal view returns (uint256) {
        uint256 num = 0;
        address[] memory oracles = IBridgeFeeOracle(oracleContract)
            .getOracleList();
        for (uint256 i = 0; i < oracles.length; i++) {
            for (uint256 j = 0; j < assets[_chainName].tokenNames.length; j++) {
                for (uint256 k = 0; k < maxQuoteIndex; k++) {
                    bytes memory asset = packAsset(
                        _chainName,
                        assets[_chainName].tokenNames[j],
                        oracles[i],
                        k
                    );
                    if (quotes[asset].expiry >= block.timestamp) {
                        num += 1;
                    }
                }
            }
        }
        return num;
    }

    function getQuote(
        string memory _chainName,
        string memory _tokenName,
        address _oracle,
        uint256 _index
    ) external view returns (QuoteInfo memory) {
        bytes memory asset = packAsset(_chainName, _tokenName, _oracle, _index);
        return
            QuoteInfo({
                id: quotes[asset].id,
                chainName: _chainName,
                tokenName: _tokenName,
                oracle: _oracle,
                fee: quotes[asset].fee,
                gasLimit: quotes[asset].gasLimit,
                expiry: quotes[asset].expiry
            });
    }

    function supportChainNames() external view returns (string[] memory) {
        return chainNames;
    }

    function supportAssets(
        string memory _chainName
    ) external view returns (Asset memory) {
        return assets[_chainName];
    }

    function verifyInput(QuoteInput memory _input) private {
        if (!assets[_input.chainName].isActive) {
            revert ChainNameInvalid();
        }
        if (!isActiveTokenName(_input.chainName, _input.tokenName)) {
            revert TokenNameInvalid();
        }
        if (
            !IBridgeFeeOracle(oracleContract).isOnline(
                _input.chainName,
                _input.oracle
            )
        ) {
            revert OracleInvalid();
        }
        if (_input.quoteIndex >= maxQuoteIndex) {
            revert QuoteIndexInvalid();
        }
        if (_input.expiry < block.timestamp) {
            revert QuoteExpired();
        }
    }

    function verifySignature(QuoteInput memory _input) private pure {
        bytes32 hash = makeMessageHash(
            _input.chainName,
            _input.tokenName,
            _input.fee,
            _input.gasLimit,
            _input.expiry
        );
        address signer = hash.toEthSignedMessageHash().recover(
            _input.signature
        );
        if (_input.oracle != signer) {
            revert VerifySignatureFailed(_input.oracle, signer);
        }
    }

    function makeMessageHash(
        string memory _chainName,
        string memory _tokenName,
        uint256 _fee,
        uint256 _gasLimit,
        uint256 _expiry
    ) public pure returns (bytes32) {
        return
            keccak256(
                abi.encode(_chainName, _tokenName, _fee, _gasLimit, _expiry)
            );
    }

    function packAsset(
        string memory _chainName,
        string memory _tokenName,
        address _oracle,
        uint256 _index
    ) internal pure returns (bytes memory) {
        return abi.encode(_chainName, _tokenName, _oracle, _index);
    }

    function unpackAsset(
        bytes memory _packedData
    )
        internal
        pure
        returns (
            string memory chainName,
            string memory tokenName,
            address oracle,
            uint256 index
        )
    {
        (chainName, tokenName, oracle, index) = abi.decode(
            _packedData,
            (string, string, address, uint256)
        );
    }

    function isActiveTokenName(
        string memory _chainName,
        string memory _tokenName
    ) public view returns (bool) {
        Asset memory asset = assets[_chainName];
        for (uint256 i = 0; i < asset.tokenNames.length; i++) {
            if (asset.tokenNames[i].equal(_tokenName)) {
                return asset.isActive;
            }
        }
        return false;
    }

    function activeTokenNames(
        string memory _chainName
    ) external view returns (string[] memory) {
        return assets[_chainName].tokenNames;
    }

    function hasActiveQuote(address _oracle) internal view returns (bool) {
        for (uint256 i = 0; i < chainNames.length; i++) {
            for (
                uint256 j = 0;
                j < assets[chainNames[i]].tokenNames.length;
                j++
            ) {
                for (uint256 k = 0; k < maxQuoteIndex; k++) {
                    bytes memory asset = packAsset(
                        chainNames[i],
                        assets[chainNames[i]].tokenNames[j],
                        _oracle,
                        k
                    );
                    if (quotes[asset].expiry >= block.timestamp) {
                        return true;
                    }
                }
            }
        }
        return false;
    }

    function registerChain(
        string memory _chainName,
        string[] memory _tokenNames
    ) external onlyOwner returns (bool) {
        if (assets[_chainName].isActive) {
            revert ChainNameAlreadyExists();
        }
        assets[_chainName] = Asset({isActive: true, tokenNames: _tokenNames});
        chainNames.push(_chainName);
        return true;
    }

    function registerTokenName(
        string memory _chainName,
        string[] memory _tokenNames
    ) external onlyOwner returns (bool) {
        if (!assets[_chainName].isActive) {
            revert ChainNameInvalid();
        }
        for (uint256 i = 0; i < _tokenNames.length; i++) {
            if (_tokenNames[i].equal("")) {
                revert TokenNameInvalid();
            }
            if (isActiveTokenName(_chainName, _tokenNames[i])) {
                revert TokenNameAlreadyExists();
            }
            assets[_chainName].tokenNames.push(_tokenNames[i]);
        }
        return true;
    }

    function updateOracleContract(
        address _oracleContract
    ) external onlyOwner returns (bool) {
        oracleContract = _oracleContract;
        return true;
    }

    function updateMaxQuoteIndex(
        uint256 _maxQuoteIndex
    ) external onlyOwner returns (bool) {
        maxQuoteIndex = _maxQuoteIndex;
        return true;
    }

    receive() external payable {}

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
