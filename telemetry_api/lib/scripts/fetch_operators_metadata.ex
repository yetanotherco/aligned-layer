defmodule Scripts.FetchOperatorsMetadata do
  require Logger
  alias TelemetryApi.ContractManagers.OperatorStateRetriever
  alias TelemetryApi.Operators.Operator
  alias TelemetryApi.Repo

  # This Script is to fetch operators metadata from the blockchain activity
  # and insert/update them into the Ecto database

  def run() do
    "Fetching old operators changes" |> Logger.debug()
    update_operators_metadata()

    "Done" |> Logger.debug()
  end

  def update_operators_metadata() do
    with {:ok, operators} <- OperatorStateRetriever.get_operators() do
      # Construct tuple {%Operator{}, op_data}
      operators = Enum.map(operators, fn op_data ->
        {Repo.get(Operator, op_data.address), op_data}
      end)

      # Fetch metadata for all operators
      operators = Enum.map(operators, fn {op, op_data} ->
        case TelemetryApi.Operators.add_operator_metadata(op_data) do
          {:ok, data} -> {:ok, {op, data}}
          {:error, msg} -> {:error, msg}
        end
      end)
        |> tap(&dbg/1)
      # Filter status ok and map to {op, op_data}
        |> Enum.filter(fn {status, _} -> status == :ok end)
        |> Enum.map(fn {_, data} -> data end)

      dbg(operators)

      # Insert in db
      Enum.map(operators, fn {op, op_data} ->
        Operator.changeset(op, op_data) |> Repo.insert_or_update()
      end)
      end
    :ok
  end
end
