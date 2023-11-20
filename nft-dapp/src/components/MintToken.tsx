import { useContractWrite, useWaitForTransaction } from "wagmi";

import { wagmiContractConfig } from "./contracts";
import { stringify } from "../utils/stringify";

export function MintToken() {
  const { write, data, error, isLoading, isError } = useContractWrite({
    ...wagmiContractConfig,
    functionName: "mint",
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
          const collection = formData.get("collection") as `0x${string}`;
          const recepient = formData.get("recepient") as `0x${string}`;
          const token = Number.parseInt(formData.get("token") as string);
          write({
            args: [collection, recepient, BigInt(token)],
          });
        }}
      >
        <input name="collection" placeholder="0xabc..." />
        <input name="recepient" placeholder="0xcde" />
        <input name="token" placeholder="1" />
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
