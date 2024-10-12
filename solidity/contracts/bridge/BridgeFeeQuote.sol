// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IERC20MetadataUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20MetadataUpgradeable.sol";
import {SafeERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/utils/SafeERC20Upgradeable.sol";
import {ECDSAUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/ECDSAUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import {IBridgeFeeQuote} from "./IBridgeFee.sol";
import {IBridgeFeeOracle} from "./IBridgeFee.sol";

error ChainNameInvalid();
error TokenInvalid();
error OracleInvalid();
error QuoteExpired();
error VerifySignatureFailed(address, address);
error QuoteNotFound();
error ChainNameAlreadyExists();
error TokenAlreadyExists();

contract BridgeFeeQuote is
    IBridgeFeeQuote,
    Initializable,
    UUPSUpgradeable,
    OwnableUpgradeable,
    ReentrancyGuardUpgradeable
{
    using SafeERC20Upgradeable for IERC20MetadataUpgradeable;
    using ECDSAUpgradeable for bytes32;

    struct Asset {
        bool isActive;
        address[] tokens;
    }

    struct Quote {
        uint256 id;
        uint256 fee;
        uint256 gasLimit;
        uint256 expiry;
    }

    address public oracleContract;

    uint256 public quoteNonce;

    string[] public chainNames;

    // chainName -> Assert
    mapping(string => Asset) public assets;

    // Only one quote is allowed per chainName + token + oracle
    mapping(bytes => Quote) internal quotes; // key: chainName + token + oracle
    // id -> chainName + token + oracle
    mapping(uint256 => bytes) internal quoteIds;

    function initialize(address _oracleContract) public initializer {
        oracleContract = _oracleContract;

        __Ownable_init();
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
    }

    event NewQuote(
        uint256 indexed id,
        address indexed oracle,
        string indexed chainName,
        address token,
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
                _inputs[i].token,
                _inputs[i].oracle
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
                _inputs[i].token,
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
            for (uint256 j = 0; j < assets[_chainName].tokens.length; j++) {
                bytes memory asset = packAsset(
                    _chainName,
                    assets[_chainName].tokens[j],
                    oracles[i]
                );
                if (quotes[asset].expiry >= block.timestamp) {
                    quotesList[currentIndex] = QuoteInfo({
                        id: quotes[asset].id,
                        chainName: _chainName,
                        token: assets[_chainName].tokens[j],
                        oracle: oracles[i],
                        fee: quotes[asset].fee,
                        gasLimit: quotes[asset].gasLimit,
                        expiry: quotes[asset].expiry
                    });
                    currentIndex += 1;
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
        (string memory chainName, address token, address oracle) = unpackAsset(
            asset
        );
        q = QuoteInfo({
            id: _id,
            chainName: chainName,
            token: token,
            oracle: oracle,
            fee: quotes[asset].fee,
            gasLimit: quotes[asset].gasLimit,
            expiry: quotes[asset].expiry
        });
    }

    /**
     * @notice Get the quote by token.
     * @param _chainName The name of the chain.
     * @param _token The address of the token.
     * @param _amount The bridge fee amount of the token.
     * @return QuoteInfo The quote.
     * @return bool Whether the quote is expired.
     */
    function getQuoteByToken(
        string memory _chainName,
        address _token,
        uint256 _amount
    ) external view returns (QuoteInfo memory, bool) {
        if (!assets[_chainName].isActive) {
            revert ChainNameInvalid();
        }
        address oracle = IBridgeFeeOracle(oracleContract).defaultOracle();

        bytes memory asset = packAsset(_chainName, _token, oracle);
        return (
            QuoteInfo({
                id: quotes[asset].id,
                chainName: _chainName,
                token: _token,
                oracle: oracle,
                fee: quotes[asset].fee,
                gasLimit: quotes[asset].gasLimit,
                expiry: quotes[asset].expiry
            }),
            quotes[asset].expiry > block.timestamp &&
                _amount >= quotes[asset].fee
        );
    }

    function currentActiveQuoteNum(
        string memory _chainName
    ) internal view returns (uint256) {
        uint256 num = 0;
        address[] memory oracles = IBridgeFeeOracle(oracleContract)
            .getOracleList();
        for (uint256 i = 0; i < oracles.length; i++) {
            for (uint256 j = 0; j < assets[_chainName].tokens.length; j++) {
                bytes memory asset = packAsset(
                    _chainName,
                    assets[_chainName].tokens[j],
                    oracles[i]
                );
                if (quotes[asset].expiry >= block.timestamp) {
                    num += 1;
                }
            }
        }
        return num;
    }

    function verifyInput(QuoteInput memory _input) private {
        if (!assets[_input.chainName].isActive) {
            revert ChainNameInvalid();
        }
        if (!isActiveToken(_input.chainName, _input.token)) {
            revert TokenInvalid();
        }
        if (
            !IBridgeFeeOracle(oracleContract).isOnline(
                _input.chainName,
                _input.oracle
            )
        ) {
            revert OracleInvalid();
        }
        if (_input.expiry < block.timestamp) {
            revert QuoteExpired();
        }
    }

    function verifySignature(QuoteInput memory _input) private pure {
        bytes32 hash = makeMessageHash(
            _input.chainName,
            _input.token,
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
        address _token,
        uint256 _fee,
        uint256 _gasLimit,
        uint256 _expiry
    ) public pure returns (bytes32) {
        return
            keccak256(abi.encode(_chainName, _token, _fee, _gasLimit, _expiry));
    }

    function packAsset(
        string memory _chainName,
        address _token,
        address _oracle
    ) internal pure returns (bytes memory) {
        return abi.encode(_chainName, _token, _oracle);
    }

    function unpackAsset(
        bytes memory _packedData
    )
        internal
        pure
        returns (string memory chainName, address token, address oracle)
    {
        (chainName, token, oracle) = abi.decode(
            _packedData,
            (string, address, address)
        );
    }

    function isActiveToken(
        string memory _chainName,
        address _token
    ) public view returns (bool) {
        Asset memory asset = assets[_chainName];
        for (uint256 i = 0; i < asset.tokens.length; i++) {
            if (asset.tokens[i] == _token) {
                return asset.isActive;
            }
        }
        return false;
    }

    function activeTokens(
        string memory _chainName
    ) external view returns (address[] memory) {
        return assets[_chainName].tokens;
    }

    function hasActiveQuote(address _oracle) internal view returns (bool) {
        for (uint256 i = 0; i < chainNames.length; i++) {
            for (uint256 j = 0; j < assets[chainNames[i]].tokens.length; j++) {
                bytes memory asset = packAsset(
                    chainNames[i],
                    assets[chainNames[i]].tokens[j],
                    _oracle
                );
                if (quotes[asset].expiry >= block.timestamp) {
                    return true;
                }
            }
        }
        return false;
    }

    function registerChain(
        string memory _chainName,
        address[] memory _tokens
    ) external onlyOwner returns (bool) {
        if (assets[_chainName].isActive) {
            revert ChainNameAlreadyExists();
        }
        assets[_chainName] = Asset({isActive: true, tokens: _tokens});
        chainNames.push(_chainName);
        return true;
    }

    function registerToken(
        string memory _chainName,
        address[] memory _tokens
    ) external onlyOwner returns (bool) {
        if (!assets[_chainName].isActive) {
            revert ChainNameInvalid();
        }
        for (uint256 i = 0; i < _tokens.length; i++) {
            if (_tokens[i] == address(0)) {
                revert TokenInvalid();
            }
            if (isActiveToken(_chainName, _tokens[i])) {
                revert TokenAlreadyExists();
            }
            assets[_chainName].tokens.push(_tokens[i]);
        }

        return true;
    }

    function updateOracleContract(
        address _oracleContract
    ) external onlyOwner returns (bool) {
        oracleContract = _oracleContract;
        return true;
    }

    receive() external payable {}

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
