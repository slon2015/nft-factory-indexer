import { configureChains, createConfig } from "wagmi";
import { goerli } from "wagmi/chains";
import { CoinbaseWalletConnector } from "wagmi/connectors/coinbaseWallet";
import { MetaMaskConnector } from "wagmi/connectors/metaMask";

import { jsonRpcProvider } from "wagmi/providers/jsonRpc";

const { chains, publicClient, webSocketPublicClient } = configureChains(
  [goerli],
  [
    jsonRpcProvider({
      rpc: (chain) =>
        chain.id == goerli.id
          ? { http: import.meta.env.VITE_GOERLY_PROVIDER as string }
          : null,
    }),
  ]
);

export const config = createConfig({
  autoConnect: true,
  connectors: [
    new MetaMaskConnector({ chains }),
    new CoinbaseWalletConnector({
      chains,
      options: {
        appName: "wagmi",
      },
    }),
  ],
  publicClient,
  webSocketPublicClient,
});
