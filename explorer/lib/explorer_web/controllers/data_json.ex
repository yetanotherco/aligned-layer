defmodule ExplorerWeb.DataJSON do
  def show(%{
        amount_of_proofs: amount_of_proofs,
        avg_fee_per_proof_usd: avg_fee_per_proof_usd
      }) do
    %{
      amount_of_proofs: amount_of_proofs,
      avg_fee_per_proof_usd: avg_fee_per_proof_usd
    }
  end
end
