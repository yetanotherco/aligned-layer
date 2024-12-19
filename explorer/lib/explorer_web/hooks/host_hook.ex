defmodule ExplorerWeb.HostHook do
  def on_mount(:add_host, _params, _session, socket) do
    %URI{host: host} = Phoenix.LiveView.get_connect_info(socket, :uri)

    socket = Phoenix.Component.assign(socket, host: host)

    {:cont, socket}
  end
end
