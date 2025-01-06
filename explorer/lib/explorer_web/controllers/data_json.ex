defmodule ExplorerWeb.DataJSON do
  def show(%{
        amount_of_proofs: amount_of_proofs,
        avg_fee_per_proof: avg_fee_per_proof
      }) do
    %{
      amount_of_proofs: amount_of_proofs,
      avg_fee_per_proof: avg_fee_per_proof
    }
  end
end
