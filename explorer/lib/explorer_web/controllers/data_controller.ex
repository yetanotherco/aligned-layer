defmodule ExplorerWeb.DataController do
  use ExplorerWeb, :controller

  def verified_proofs_in_last_24_hours(conn, _params) do
    %{
      amount_of_proofs: amount_of_proofs,
      avg_fee_per_proof: avg_fee_per_proof
    } = Batches.get_verified_proofs_in_last_24_hours()

    render(conn, :show, %{
      amount_of_proofs: amount_of_proofs,
      avg_fee_per_proof: avg_fee_per_proof
    })
  end
end
