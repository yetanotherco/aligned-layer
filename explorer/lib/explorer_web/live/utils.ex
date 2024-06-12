defmodule ExplorerWeb.Utils do
  def shorten_hash(hash) do
    case String.length(hash) do
      n when n < 6 -> hash
      _ -> "#{String.slice(hash, 0, 6)}...#{String.slice(hash, -4, 4)}"
    end
  end

  def convert_number_to_shorthand(number) when number >= 1_000_000 do
    "#{div(number, 1_000_000)}M"
  end

  def convert_number_to_shorthand(number) when number >= 10_000 do
    "#{div(number, 10_000)}k"
  end

  def convert_number_to_shorthand(number) when number >= 1_000 do
    "#{div(number, 1_000)}k"
  end

  def convert_number_to_shorthand(number) when number >= 0 do
    "#{number}"
  end

  def convert_number_to_shorthand(_number), do: "Invalid number"
end

defmodule Utils do
  def string_to_bytes32(hex_string) do
    # Remove the '0x' prefix
    hex = case hex_string do
      "0x" <> _ -> String.slice(hex_string, 2..-1//1)
      _ -> raise "Invalid hex string, missing '0x' prefix"
    end

    # Convert the hex string to a binary
    case Base.decode16(hex, case: :mixed) do
      {:ok, binary} -> binary
      _ -> raise "Invalid hex string"
    end
  end

  def get_last_n_items(events, n) when is_list(events) and is_integer(n) and n >= 0 do
    events
    |> Enum.reverse()
    |> Enum.take(n)
    |> Enum.reverse()
  end

  def extract_amount_of_proofs_from_json({:ok, batch_json}) do
    batch_json |> Enum.count()
  end

  def extract_amount_of_proofs_from_json({:error, _}) do
    300
  end

  def fetch_batch_data_pointer(batch_data_pointer) do # TODO fix Bottleneck
    case Finch.build(:get, batch_data_pointer) |> Finch.request(Explorer.Finch) do
      {:ok, %Finch.Response{status: 200, body: body}} ->
        case Jason.decode(body) do
          {:ok, json} -> {:ok, json}
          {:error, reason} -> {:error, {:json_decode, reason}}
        end
      {:ok, %Finch.Response{status: status_code}} ->
        {:error, {:http_error, status_code}}
      {:error, reason} ->
        {:error, {:http_error, reason}}
    end
  end

  def extract_amount_of_proofs(%BatchDB{} = batch) do
    amount_of_proofs = batch.data_pointer |> Utils.fetch_batch_data_pointer |> Utils.extract_amount_of_proofs_from_json
    Map.put(batch, :amount_of_proofs, amount_of_proofs)
  end
end
