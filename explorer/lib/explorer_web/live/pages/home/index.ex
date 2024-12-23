defmodule ExplorerWeb.Home.Index do
  require Logger
  import ExplorerWeb.ChartComponents
  import ExplorerWeb.BatchesTable
  use ExplorerWeb, :live_view

  def get_stats() do
    verified_batches = Batches.get_amount_of_verified_batches()
    avg_fee_per_proof = Batches.get_avg_fee_per_proof()

    avg_fee_per_proof_usd =
      case EthConverter.wei_to_usd(avg_fee_per_proof, 2) do
        {:ok, value} -> value
        _ -> 0
      end

    avg_fee_per_proof_eth = EthConverter.wei_to_eth(avg_fee_per_proof, 4)

    verified_proofs = Batches.get_amount_of_verified_proofs()
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

  def get_cost_per_proof_chart_data() do
    batches = Enum.reverse(Batches.get_latest_batches(%{amount: 100, order_by: :desc}))

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
    batches = Enum.reverse(Batches.get_latest_batches(%{amount: 100, order_by: :desc}))

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

  defp set_empty_values(socket) do
    Logger.info("Setting empty values")

    socket
    |> assign(
      stats: [],
      latest_batches: [],
      cost_per_proof_data: %{points: [], extra_data: %{}},
      batch_size_chart_data: %{points: [], extra_data: %{}}
    )
  end

  @impl true
  def handle_info(_, socket) do
    latest_batches = Batches.get_latest_batches(%{amount: 5, order_by: :desc})

    {:noreply,
     assign(
       socket,
       stats: get_stats(),
       latest_batches: latest_batches,
       cost_per_proof_chart: get_cost_per_proof_chart_data(),
       batch_size_chart_data: get_batch_size_chart_data()
     )}
  end

  @impl true
  def mount(_, _, socket) do
    latest_batches = Batches.get_latest_batches(%{amount: 10, order_by: :desc})

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
      Logger.error("Error: Mint.TransportError: #{inspect(e)}")

      case e do
        %Mint.TransportError{reason: :econnrefused} ->
          {
            :ok,
            set_empty_values(socket)
            |> put_flash(:error, "Could not connect to the backend, please try again later.")
          }

        _ ->
          {
            :ok,
            set_empty_values(socket)
            |> put_flash(:error, "Something went wrong, please try again later.")
          }
      end

    e in FunctionClauseError ->
      Logger.error("Error: FunctionClauseError: #{inspect(e)}")

      case e do
        %FunctionClauseError{
          module: ExplorerWeb.Home.Index
        } ->
          {
            :ok,
            set_empty_values(socket)
            |> put_flash(:error, "Something went wrong with the RPC, please try again later.")
          }
      end

    e ->
      Logger.error("Error: other error: #{inspect(e)}")

      {
        :ok,
        set_empty_values(socket)
        |> put_flash(:error, "Something went wrong, please try again later.")
      }
  end

  embed_templates("*")
end
