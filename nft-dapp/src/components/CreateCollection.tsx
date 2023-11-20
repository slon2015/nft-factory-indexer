import { useContractWrite, useWaitForTransaction } from "wagmi";

import { wagmiContractConfig } from "./contracts";
import { stringify } from "../utils/stringify";

export function CreateCollection() {
  const { write, data, error, isLoading, isError } = useContractWrite({
    ...wagmiContractConfig,
    functionName: "createCollection",
  });
  const {
    data: receipt,
    isLoading: isPending,
    isSuccess,
  } = useWaitForTransaction({ hash: data?.hash });

  return (
    <>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          const formData = new FormData(e.target as HTMLFormElement);
          const name = formData.get("name") as string;
          const symbol = formData.get("symbol") as string;
          const uri = formData.get("uri") as string;
          write({
            args: [name, symbol, uri],
          });
        }}
      >
        <input name="name" placeholder="name" />
        <input name="symbol" placeholder="symbol" />
        <input name="uri" placeholder="http://smtn.com" />
        <button type="submit">Send</button>
      </form>

      {isLoading && <div>Check wallet...</div>}
      {isPending && <div>Transaction pending...</div>}
      {isSuccess && (
        <>
          <div>Transaction Hash: {data?.hash}</div>
          <div>
            Transaction Receipt: <pre>{stringify(receipt, null, 2)}</pre>
          </div>
        </>
      )}
      {isError && <div>Error: {error?.message}</div>}
    </>
  );
}
