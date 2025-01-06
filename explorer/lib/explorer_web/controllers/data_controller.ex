defmodule ExplorerWeb.DataController do
  use ExplorerWeb, :controller

  def last_verified_proofs_count(conn, _params) do
    last_verified_proofs_count = Batches.get_last_verified_proofs_count()
    render(conn, :show, count: last_verified_proofs_count)
  end
end
