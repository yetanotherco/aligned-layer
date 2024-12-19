defmodule ExplorerWeb.Home.Index do
  require Logger
  import ExplorerWeb.ChartComponents
  use ExplorerWeb, :live_view

  def get_cost_per_proof_chart_data() do
    batches = Batches.get_latest_batches(%{amount: 100, order_by: :asc})

    extra_data =
      %{
        merkle_root: Enum.map(batches, fn b -> b.merkle_root end),
        amount_of_proofs: Enum.map(batches, fn b -> b.amount_of_proofs end),
        age: Enum.map(batches, fn b -> Helpers.parse_timeago(b.submission_timestamp) end)
      }

    points =
      Enum.map(batches, fn b ->
        fee_per_proof =
          case EthConverter.wei_to_usd(b.fee_per_proof, 2) do
            {:ok, value} ->
              value

            # Nil values are ignored by the chart
            {:error, _} ->
              nil
          end

        %{x: b.submission_block_number, y: fee_per_proof}
      end)

    %{
      points: points,
      extra_data: extra_data
    }
  end

  @impl true
  def handle_info(_, socket) do
    verified_batches = Batches.get_amount_of_verified_batches()

    operators_registered = Operators.get_amount_of_operators()

    latest_batches =
      Batches.get_latest_batches(%{amount: 5, order_by: :desc})
      # extract only the merkle root
      |> Enum.map(fn %Batches{merkle_root: merkle_root} -> merkle_root end)

    verified_proofs = Batches.get_amount_of_verified_proofs()

    restaked_amount_eth = Restakings.get_restaked_amount_eth()
    restaked_amount_usd = Restakings.get_restaked_amount_usd()

    {:noreply,
     assign(
       socket,
       verified_batches: verified_batches,
       operators_registered: operators_registered,
       latest_batches: latest_batches,
       verified_proofs: verified_proofs,
       restaked_amount_eth: restaked_amount_eth,
       restaked_amount_usd: restaked_amount_usd,
       cost_per_proof_chart: get_cost_per_proof_chart_data()
     )}
  end

  @impl true
  def mount(_, _, socket) do
    verified_batches = Batches.get_amount_of_verified_batches()

    operators_registered = Operators.get_amount_of_operators()

    latest_batches =
      Batches.get_latest_batches(%{amount: 5, order_by: :desc})
      # extract only the merkle root
      |> Enum.map(fn %Batches{merkle_root: merkle_root} -> merkle_root end)

    verified_proofs = Batches.get_amount_of_verified_proofs()

    restaked_amount_eth = Restakings.get_restaked_amount_eth()
    restaked_amount_usd = Restakings.get_restaked_amount_usd()

    if connected?(socket), do: Phoenix.PubSub.subscribe(Explorer.PubSub, "update_views")

    {:ok,
     assign(socket,
       verified_batches: verified_batches,
       operators_registered: operators_registered,
       latest_batches: latest_batches,
       verified_proofs: verified_proofs,
       service_manager_address:
         AlignedLayerServiceManager.get_aligned_layer_service_manager_address(),
       restaked_amount_eth: restaked_amount_eth,
       restaked_amount_usd: restaked_amount_usd,
       cost_per_proof_chart: get_cost_per_proof_chart_data(),
       page_title: "Welcome"
     )}
  rescue
    e in Mint.TransportError ->
      case e do
        %Mint.TransportError{reason: :econnrefused} ->
          {
            :ok,
            assign(socket,
              verified_batches: :empty,
              operators_registered: :empty,
              latest_batches: :empty,
              verified_proofs: :empty
            )
            |> put_flash(:error, "Could not connect to the backend, please try again later.")
          }

        _ ->
          "Other transport error: #{inspect(e)}" |> Logger.error()
          {:ok, socket |> put_flash(:error, "Something went wrong, please try again later.")}
      end

    e in FunctionClauseError ->
      case e do
        %FunctionClauseError{
          module: ExplorerWeb.Home.Index
        } ->
          {
            :ok,
            assign(socket,
              verified_batches: :empty,
              operators_registered: :empty,
              latest_batches: :empty,
              verified_proofs: :empty
            )
            |> put_flash(:error, "Something went wrong with the RPC, please try again later.")
          }
      end

    e ->
      Logger.error("Other error: #{inspect(e)}")
      {:ok, socket |> put_flash(:error, "Something went wrong, please try again later.")}
  end

  embed_templates("*")
end
