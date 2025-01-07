defmodule ExplorerWeb.DataController do
  use ExplorerWeb, :controller

  def verified_batches_summary(conn, _params) do
    batches = Batches.get_verified_batches_summary()

    render(conn, :show, %{
      batches: batches,
    })
  end
end
