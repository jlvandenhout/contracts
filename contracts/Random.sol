// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@iota/iscmagic/ISC.sol";

contract Random {
    event Int(uint256);
    event Bytes(bytes);

    uint256 private _nonce;

    function getNonce() internal returns (bytes32) {
        return bytes32(_nonce++);
    }

    function getInt() public returns (uint256) {
        bytes32 entropy = ISC.sandbox.getEntropy();
        bytes32 nonce = getNonce();
        bytes32 digest = keccak256(bytes.concat(entropy, nonce));

        uint256 value = uint256(digest);

        emit Int(value);
        return value;
    }

    function getBytes(uint length) public returns (bytes memory) {
        bytes32 entropy = ISC.sandbox.getEntropy();
        bytes32 nonce = getNonce();
        bytes32 digest = keccak256(bytes.concat(entropy, nonce));

        bytes memory value = new bytes(length);
        for (uint i = 0; i < length; i += 32) {
            digest = keccak256(bytes.concat(digest));
            for (uint j = 0; j < 32 && i + j < length; j++) {
                value[i + j] = digest[j];
            }
        }

        emit Bytes(value);
        return value;
    }
}
