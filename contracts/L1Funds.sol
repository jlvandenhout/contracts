// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@iota/iscmagic/ISC.sol";

contract L1Funds {
  event Send(address sender, L1Address receiver, ISCAssets assets);

  function send(L1Address memory receiver, ISCAssets memory allowance) public {
      ISC.sandbox.allow(address(this), allowance);
      
      ISCSendMetadata memory metadata;
      ISCSendOptions memory options;
      ISC.sandbox.send(receiver, allowance, true, metadata, options);

      emit Send(msg.sender, receiver, allowance);
  }
}
