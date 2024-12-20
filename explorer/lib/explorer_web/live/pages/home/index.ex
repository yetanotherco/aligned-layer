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

  def get_batch_size_chart_data() do
    batches = Batches.get_latest_batches(%{amount: 100, order_by: :asc})

    extra_data =
      %{
        merkle_root: Enum.map(batches, fn b -> b.merkle_root end),
        amount_of_proofs: Enum.map(batches, fn b -> b.amount_of_proofs end),
        age: Enum.map(batches, fn b -> Helpers.parse_timeago(b.submission_timestamp) end)
      }

    points =
      Enum.map(batches, fn b ->
        %{x: b.submission_block_number, y: b.amount_of_proofs}
      end)

    %{
      points: points,
      extra_data: extra_data
    }
  end

  def get_stats() do
    verified_proofs = Batches.get_amount_of_verified_proofs()
    verified_batches = Batches.get_amount_of_verified_batches()
    avg_fee_per_proof = Batches.get_avg_fee_per_proof()

    avg_fee_per_proof_usd =
      case EthConverter.wei_to_usd(avg_fee_per_proof, 2) do
        {:ok, value} -> value
        _ -> 0
      end

    avg_fee_per_proof_eth = EthConverter.wei_to_eth(avg_fee_per_proof, 4)

    operators_registered = Operators.get_amount_of_operators()

    restaked_amount_eth = Restakings.get_restaked_amount_eth()
    restaked_amount_usd = Restakings.get_restaked_amount_usd()
    operator_latest_release = ReleasesHelper.get_latest_release()

    [
      %{
        title: "Proofs verified",
        value: Helpers.convert_number_to_shorthand(verified_proofs),
        tooltip_text: nil
      },
      %{
        title: "Total batches",
        value: Helpers.convert_number_to_shorthand(verified_batches),
        tooltip_text: nil
      },
      %{
        title: "AVG proof cost",
        value: "#{avg_fee_per_proof_usd} USD",
        tooltip_text: "~= #{avg_fee_per_proof_eth} ETH"
      },
      %{
        title: "Operators",
        value: operators_registered,
        tooltip_text: "Current version #{operator_latest_release}"
      },
      %{
        title: "Total restake",
        value: "#{restaked_amount_usd} USD",
        tooltip_text: "~= #{restaked_amount_eth} ETH"
      }
    ]
  end

  @impl true
  def handle_info(_, socket) do
    latest_batches =
      Batches.get_latest_batches(%{amount: 10, order_by: :desc})

    {:noreply,
     assign(
       socket,
       latest_batches: latest_batches,
       stats: get_stats(),
       cost_per_proof_chart: get_cost_per_proof_chart_data(),
       batch_size_chart_data: get_batch_size_chart_data()
     )}
  end

  @impl true
  def mount(_, _, socket) do
    latest_batches =
      Batches.get_latest_batches(%{amount: 10, order_by: :desc})

    if connected?(socket), do: Phoenix.PubSub.subscribe(Explorer.PubSub, "update_views")

    {:ok,
     assign(socket,
       stats: get_stats(),
       latest_batches: latest_batches,
       cost_per_proof_chart: get_cost_per_proof_chart_data(),
       batch_size_chart_data: get_batch_size_chart_data(),
       page_title: "Welcome"
     )}
  rescue
    e in Mint.TransportError ->
      case e do
        %Mint.TransportError{reason: :econnrefused} ->
          {
            :ok,
            assign(socket,
              stats: :empty,
              latest_batches: :empty
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
