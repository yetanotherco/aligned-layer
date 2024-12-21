defmodule ExplorerWeb.Home.Index do
  require Logger
  use ExplorerWeb, :live_view

  defp set_empty_values(socket) do
    Logger.info("Setting empty values")
    socket |> assign(
      verified_batches: :empty,
      operators_registered: :empty,
      latest_batches: :empty,
      verified_proofs: :empty,
      restaked_amount_eth: :empty,
      restaked_amount_usd: :empty
    )
  end

  @impl true
  def handle_info(_, socket) do
    verified_batches = Batches.get_amount_of_verified_batches()

    operators_registered = Operators.get_amount_of_operators()

    latest_batches =
      Batches.get_latest_batches(%{amount: 5})
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
       restaked_amount_usd: restaked_amount_usd
     )}
  end

  @impl true
  def mount(_, _, socket) do
    verified_batches = Batches.get_amount_of_verified_batches()

    operators_registered = Operators.get_amount_of_operators()

    latest_batches =
      Batches.get_latest_batches(%{amount: 5})
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
            |> put_flash(:error, "Something went wrong, please try again later.")}
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
        |> put_flash(:error, "Something went wrong, please try again later.")}
  end

  embed_templates("*")
end
