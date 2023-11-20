import { useNetwork, useSwitchNetwork } from "wagmi";

export function NetworkSwitcher() {
  const { chain } = useNetwork();
  const { chains, error, isLoading, pendingChainId, switchNetwork } =
    useSwitchNetwork();

  const ableToSwitch = chains.length > 1 || chain?.unsupported;

  return (
    <div>
      <div>
        Connected to {chain?.name ?? chain?.id}
        {chain?.unsupported && " (unsupported)"}
      </div>
      {switchNetwork && ableToSwitch && (
        <>
          <br />
          <div>
            Switch to:{" "}
            {chains.map((x) =>
              x.id === chain?.id ? null : (
                <button key={x.id} onClick={() => switchNetwork(x.id)}>
                  {x.name}
                  {isLoading && x.id === pendingChainId && " (switching)"}
                </button>
              )
            )}
          </div>
        </>
      )}

      <div>{error?.message}</div>
    </div>
  );
}
