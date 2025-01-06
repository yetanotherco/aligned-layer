defmodule ExplorerWeb.DataController do
  use ExplorerWeb, :controller

  def verified_proofs_in_last_24_hours(conn, _params) do
    %{
      amount_of_proofs: amount_of_proofs,
      avg_fee_per_proof: avg_fee_per_proof
    } = Batches.get_last_24h_verified_proof_stats()

    avg_fee_per_proof_usd =
      case EthConverter.wei_to_usd_sf(avg_fee_per_proof, 2) do
        {:ok, value} -> value
        _ -> 0
      end

    render(conn, :show, %{
      amount_of_proofs: amount_of_proofs,
      avg_fee_per_proof_usd: avg_fee_per_proof_usd
    })
  end
end
