defmodule ExplorerWeb.DataJSON do
  def show(%{count: last_verified_proofs_count}) do
    %{
      count: last_verified_proofs_count
    }
  end
end
