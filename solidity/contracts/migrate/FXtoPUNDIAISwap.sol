// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Pausable} from "@openzeppelin/contracts/security/Pausable.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/security/ReentrancyGuard.sol";

interface IPundiAIToken {
    function mint(address to, uint256 amount) external;
}

contract FXtoPUNDIAISwap is Pausable, Ownable, ReentrancyGuard {
    IERC20 public immutable FX_TOKEN;
    IPundiAIToken public immutable PUNDIAI_TOKEN;

    uint256 public constant EXCHANGE_RATE = 100; // FX -> PUNDIAI rate divisor
    uint256 public immutable MAX_TOTAL_MINT; // Max total mintable amount

    uint256 public totalMinted; // Total PUNDIAI minted so far

    mapping(address => bool) public hasSwapped; // Tracks if an address has swapped
    mapping(address => bool) public blacklist; // Blacklisted addresses

    event Swap(
        address indexed from,
        address indexed to,
        uint256 fxAmount,
        uint256 pundiAmount
    );
    event BlacklistUpdated(address indexed user, bool isBlacklisted);

    constructor(address _fxToken, address _pundiAIToken) {
        require(
            _fxToken != address(0) && _pundiAIToken != address(0),
            "Invalid token address"
        );

        FX_TOKEN = IERC20(_fxToken);
        PUNDIAI_TOKEN = IPundiAIToken(_pundiAIToken);
        uint256 balance = FX_TOKEN.totalSupply() -
            FX_TOKEN.balanceOf(0xB5A58db25eeDEfeEDe888f55D6157e13d4b4f4F8);
        MAX_TOTAL_MINT = balance / EXCHANGE_RATE;
        blacklist[0xB5A58db25eeDEfeEDe888f55D6157e13d4b4f4F8] = true;
    }

    function swap() external nonReentrant whenNotPaused {
        _swap(msg.sender, msg.sender);
    }

    function swapFor(address _contract, address _to) external onlyOwner {
        require(_to != address(0), "Invalid address");
        require(isContract(_contract), "Not a contract address");
        require(!blacklist[_to], "Address is blacklisted");
        _swap(_contract, _to);
    }

    function _swap(address _from, address _to) internal {
        require(!blacklist[_from], "Address is blacklisted");
        require(!hasSwapped[_from], "Already swapped");

        uint256 fxBalance = FX_TOKEN.balanceOf(_from);
        require(fxBalance > 0, "No FX tokens to swap");

        uint256 pundiAIAmount = fxBalance / EXCHANGE_RATE;
        require(pundiAIAmount > 0, "Swap amount too small");
        require(
            totalMinted + pundiAIAmount <= MAX_TOTAL_MINT,
            "Exceeds total mint limit"
        );

        hasSwapped[_from] = true;
        totalMinted += pundiAIAmount;

        PUNDIAI_TOKEN.mint(_to, pundiAIAmount);

        emit Swap(_from, _to, fxBalance, pundiAIAmount);
    }

    function setBlacklist(address user, bool status) external onlyOwner {
        blacklist[user] = status;
        emit BlacklistUpdated(user, status);
    }

    function withdrawToken(address token, uint256 amount) external onlyOwner {
        IERC20(token).transfer(owner(), amount);
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }

    function isContract(address account) internal view returns (bool) {
        return account.code.length > 0;
    }
}
