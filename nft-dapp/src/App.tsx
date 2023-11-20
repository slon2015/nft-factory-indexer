import { useAccount } from "wagmi";

import { Account } from "./components/Account";
import { Connect } from "./components/Connect";
import { NetworkSwitcher } from "./components/NetworkSwitcher";
import { CreateCollection } from "./components/CreateCollection";
import { MintToken } from "./components/MintToken";

export function App() {
  const { isConnected } = useAccount();

  return (
    <>
      <h1>NFT Factory Admin panel</h1>

      <Connect />

      {isConnected && (
        <>
          <hr />
          <h2>Network</h2>
          <NetworkSwitcher />
          <br />
          <hr />
          <h2>Account</h2>
          <Account />
          <br />
          <hr />
          <h2>Create collection</h2>
          <CreateCollection />
          <br />
          <hr />
          <h2>Mint token</h2>
          <MintToken />
        </>
      )}
    </>
  );
}
