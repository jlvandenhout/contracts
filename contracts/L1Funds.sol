// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@iota/iscmagic/ISC.sol";

contract L1Funds {
  event BaseTokenEvent(uint64 amount);
  event NFTEvent(NFTID id);

  function deposit(ISCAssets memory assets) public {
    emit BaseTokenEvent(assets.baseTokens);
    emit NFTEvent(assets.nfts[0]);
  }
}
