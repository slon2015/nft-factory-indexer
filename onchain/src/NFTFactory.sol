// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

import {NFT} from "./NFT.sol";

contract NFTFactory is Ownable, ReentrancyGuard {

    event CollectionCreated(address collection, string name, string symbol);
    event TokenMinted(address collection, address recipient, uint256 tokenId, string tokenURI);

    mapping (address => bool) collections;

    constructor() Ownable(_msgSender()) ReentrancyGuard() {}

    function getSalt(string memory name, string memory symbol, string memory uri) internal pure returns(bytes32) {
        bytes memory content = abi.encodePacked(name, symbol, uri);
        return keccak256(content);
    }

    function createCollection(string memory name, string memory symbol, string memory uri) external onlyOwner returns (address) {
        NFT deployedCollection = new NFT{salt: getSalt(name, symbol, uri)}(name, symbol, uri);

        address deployedCollectionAddress = address(deployedCollection);

        collections[deployedCollectionAddress] = true;

        emit CollectionCreated(deployedCollectionAddress, name, symbol);

        return deployedCollectionAddress;
    }

    function mint(address collection, address to, uint256 tokenId) external onlyOwner nonReentrant {
        require(collections[collection], "Not managed contract");
        NFT deployedCollection = NFT(collection);

        deployedCollection.mint(to, tokenId);

        emit TokenMinted(collection, to, tokenId, deployedCollection.tokenURI(tokenId));
    }
}