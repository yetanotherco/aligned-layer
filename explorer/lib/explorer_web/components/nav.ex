defmodule NavComponent do
  use ExplorerWeb, :live_component

  def get_current_network(host) do
    case host do
      "explorer.alignedlayer.com" -> "Mainnet"
      "holesky.explorer.alignedlayer.com" -> "Testnet"
      "stage.explorer.alignedlayer.com" -> "Stage"
      _ -> "Devnet"
    end
  end

  @impl true
  def mount(socket) do
    networks = ExplorerWeb.Helpers.get_aligned_networks()

    networks =
      Enum.map(networks, fn {name, link} ->
        {name, "window.location.href='#{link}'"}
      end)

    {:ok,
     assign(socket,
       latest_release: ReleasesHelper.get_latest_release(),
       networks: networks
     )}
  end

  @impl true
  def render(assigns) do
    ~H"""
    <nav class={
      classes([
        "px-4 sm:px-6 lg:px-8 fixed top-0 p-3 z-50",
        "flex justify-between items-center w-full",
        "border-b border-foreground/10 backdrop-blur-lg backdrop-saturate-200"
      ])
    }>
      <div class="gap-x-6 inline-flex">
        <.link
          class="hover:scale-105 transform duration-150 active:scale-95 text-3xl"
          navigate={~p"/"}
        >
          🟩 <span class="sr-only">Aligned Explorer Home</span>
        </.link>
        <div class={["items-center gap-8 [&>a]:drop-shadow-md", "hidden md:inline-flex"]}>
          <.link
            class={
              active_view_class(assigns.socket.view, [
                ExplorerWeb.Batches.Index,
                ExplorerWeb.Batch.Index
              ])
            }
            navigate={~p"/batches"}
          >
            Batches
          </.link>
          <.link
            class={
              active_view_class(assigns.socket.view, [
                ExplorerWeb.Operators.Index,
                ExplorerWeb.Operator.Index
              ])
            }
            navigate={~p"/operators"}
          >
            Operators
          </.link>
          <.link
            class={
                active_view_class(assigns.socket.view, [
                  ExplorerWeb.Restakes.Index,
                  ExplorerWeb.Restake.Index
                ])
              }
            navigate={~p"/restakes"}
          >
            Restakes
          </.link>
        </div>
        <.live_component module={SearchComponent} id="nav_search" />
      </div>
      <div class="items-center gap-4 font-semibold leading-6 text-foreground/80 flex [&>a]:hidden lg:[&>a]:inline-block [&>a]:drop-shadow-md">
        <.link class="hover:text-foreground" target="_blank" href="https://docs.alignedlayer.com">
          Docs
        </.link>
        <.link
          class="hover:text-foreground"
          target="_blank"
          href="https://github.com/yetanotherco/aligned_layer"
        >
          GitHub
        </.link>
        <DarkMode.button />
        <.badge :if={@latest_release != nil} class="hidden md:inline">
          <%= @latest_release %>
          <.tooltip>
            Latest Aligned version
          </.tooltip>
        </.badge>
        <.hover_dropdown_selector
          current_value={get_current_network(@host)}
          variant="foreground"
          options={@networks}
          icon="hero-cube-transparent-micro"
        />
        <button
          class="md:hidden z-50"
          id="menu-toggle"
          phx-click={toggle_menu()}
          aria-label="Toggle hamburger menu"
        >
          <.icon name="hero-bars-3" class="toggle-open" />
          <.icon name="hero-x-mark" class="toggle-close hidden" />
        </button>
        <div
          id="menu-overlay"
          class="fixed inset-0 bg-background/90 z-40 hidden min-h-dvh animate-in fade-in"
          phx-click={toggle_menu()}
        >
          <div class="h-full flex flex-col gap-y-10 text-2xl justify-end items-center p-12">
            <.badge :if={@latest_release != nil}>
              <%= @latest_release %>
            </.badge>
            <.link
              class={
                classes([
                  active_view_class(assigns.socket.view, [
                    ExplorerWeb.Batches.Index,
                    ExplorerWeb.Batch.Index
                  ]),
                  "text-foreground/80 hover:text-foreground font-semibold"
                ])
              }
              navigate={~p"/batches"}
            >
              Batches
            </.link>
            <.link
              class={
                active_view_class(assigns.socket.view, [
                  ExplorerWeb.Operators.Index,
                  ExplorerWeb.Operator.Index
                ])
              }
              navigate={~p"/operators"}
            >
              Operators
            </.link>
            <.link
              class={
                active_view_class(assigns.socket.view, [
                  ExplorerWeb.Restakes.Index,
                  ExplorerWeb.Restake.Index
                ])
              }
              navigate={~p"/restakes"}
            >
              Restakes
            </.link>
            <.link class="hover:text-foreground" target="_blank" href="https://docs.alignedlayer.com">
              Docs
            </.link>
            <.link
              class="hover:text-foreground"
              target="_blank"
              href="https://github.com/yetanotherco/aligned_layer"
            >
              GitHub
            </.link>
          </div>
        </div>
      </div>
    </nav>
    """
  end

  def toggle_menu() do
    JS.toggle(to: "#menu-overlay")
    |> JS.toggle(to: ".toggle-open")
    |> JS.toggle(to: ".toggle-close")
  end

  defp active_view_class(current_view, target_views) do
    if current_view in target_views,
      do: "text-green-500 font-bold",
      else: "text-foreground/80 hover:text-foreground font-semibold"
  end
end
