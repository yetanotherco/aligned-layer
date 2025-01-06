defmodule ExplorerWeb.DataController do
  use ExplorerWeb, :controller

  def verified_proofs_in_last_24_hours(conn, _params) do
    verified_proofs_in_last_24_hours = Batches.get_verified_proofs_in_last_24_hours()
    render(conn, :show, count: verified_proofs_in_last_24_hours)
  end
end
