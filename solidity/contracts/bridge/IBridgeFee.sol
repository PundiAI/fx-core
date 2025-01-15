// SPDX-License-Identifier: UNLICENSED
/* solhint-disable one-contract-per-file */
pragma solidity ^0.8.0;

interface IBridgeFeeQuote {
    struct Quote {
        uint64 expiry;
        uint64 gasLimit;
        uint256 amount;
        bytes32 index;
    }

    struct QuoteIndex {
        uint256 id;
        bytes32 chainName;
        bytes32 tokenName;
        address oracle;
        uint8 cap;
    }

    struct QuoteInput {
        uint8 cap;
        uint64 gasLimit;
        uint64 expiry;
        bytes32 chainName;
        bytes32 tokenName;
        uint256 amount;
    }

    struct QuoteInfo {
        uint256 id;
        bytes32 chainName;
        bytes32 tokenName;
        address oracle;
        uint256 amount;
        uint64 gasLimit;
        uint64 expiry;
    }

    function quote(
        QuoteInput[] calldata _inputs
    ) external returns (uint256[] memory);

    function getQuoteById(uint256 _id) external view returns (QuoteInfo memory);

    function getQuotesByToken(
        bytes32 _chainName,
        bytes32 _token
    ) external view returns (QuoteInfo[] memory quotes);

    function getQuoteByIndex(
        bytes32 _chainName,
        bytes32 _token,
        address _oracle,
        uint8 _cap
    ) external view returns (QuoteInfo memory);

    function getDefaultOracleQuote(
        bytes32 _chainName,
        bytes32 _token
    ) external view returns (QuoteInfo[] memory quotes);

    function getChainNames() external view returns (bytes32[] memory);

    function getTokens(
        bytes32 _chainName
    ) external view returns (bytes32[] memory);

    event NewQuote(
        uint256 indexed id,
        bytes32 indexed chainName,
        bytes32 indexed tokenName,
        address oracle,
        uint256 fee,
        uint256 gasLimit,
        uint256 expiry,
        uint8 cap
    );
}

interface IBridgeFeeOracle {
    function defaultOracle() external view returns (address);

    function isOnline(
        bytes32 _chainName,
        address _oracle
    ) external returns (bool);

    function getOracleList(
        bytes32 _chainName
    ) external view returns (address[] memory);
}
