// SPDX-License-Identifier: UNLICENSED
/* solhint-disable one-contract-per-file */
pragma solidity ^0.8.0;

interface IBridgeFeeQuote {
    struct Asset {
        bool isActive;
        string[] tokenNames;
    }

    struct QuoteInput {
        string chainName;
        string tokenName;
        address oracle;
        uint256 quoteIndex;
        uint256 fee;
        uint256 gasLimit;
        uint256 expiry;
        bytes signature;
    }

    struct QuoteInfo {
        uint256 id;
        string chainName;
        string tokenName;
        address oracle;
        uint256 fee;
        uint256 gasLimit;
        uint256 expiry;
    }

    function quote(QuoteInput[] memory _inputs) external returns (bool);

    function getQuoteList(
        string memory _chainName
    ) external view returns (QuoteInfo[] memory);

    function getQuoteById(uint256 _id) external view returns (QuoteInfo memory);

    function getQuotesByToken(
        string memory _chainName,
        string memory _tokenName
    ) external view returns (QuoteInfo[] memory quotes);

    function getQuote(
        string memory _chainName,
        string memory _tokenName,
        address _oracle,
        uint256 _index
    ) external view returns (QuoteInfo memory);

    function supportChainNames() external view returns (string[] memory);

    function supportAssets(
        string memory _chainName
    ) external view returns (Asset memory);

    function maxQuoteIndex() external view returns (uint256);
}

interface IBridgeFeeOracle {
    function defaultOracle() external view returns (address);

    function isOnline(
        string memory _chainName,
        address _oracle
    ) external returns (bool);

    function getOracleList() external view returns (address[] memory);
}
