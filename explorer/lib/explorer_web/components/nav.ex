defmodule NavComponent do
  use ExplorerWeb, :live_component

  @impl true
  def mount(socket) do
    {:ok, assign(socket, latest_release: ReleasesHelper.get_latest_release())}
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
          class={
            classes([
              "hover:scale-105",
              "transform",
              "duration-150",
              "active:scale-95",
              active_view_class(assigns.socket.view, ExplorerWeb.Home.Index)
            ])
          }
          navigate={~p"/"}
        >
          🟩 <span class="sr-only">Aligned Explorer Home</span>
        </.link>
        <div class={["items-center gap-8 [&>a]:drop-shadow-md", "hidden md:inline-flex"]}>
          <.link
            class={active_view_class(assigns.socket.view, ExplorerWeb.Batches.Index)}
            navigate={~p"/batches"}
          >
            Batches
          </.link>
          <.link
            class={active_view_class(assigns.socket.view, ExplorerWeb.Operators.Index)}
            navigate={~p"/operators"}
          >
            Operators
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
                  active_view_class(assigns.socket.view, ExplorerWeb.Batches.Index),
                  "text-foreground/80 hover:text-foreground font-semibold"
                ])
              }
              navigate={~p"/batches"}
            >
              Batches
            </.link>
            <.link
              class={active_view_class(assigns.socket.view, ExplorerWeb.Operators.Index)}
              navigate={~p"/operators"}
            >
              Operators
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

  defp active_view_class(ExplorerWeb.Home.Index, ExplorerWeb.Home.Index), do: "text-4xl"
  defp active_view_class(_other, ExplorerWeb.Home.Index), do: "text-3xl"

  defp active_view_class(ExplorerWeb.Batches.Index, ExplorerWeb.Batches.Index ), do: "text-green-500 font-bold"
  defp active_view_class(ExplorerWeb.Batch.Index, ExplorerWeb.Batches.Index ), do: "text-green-500 font-bold"
  defp active_view_class(_other, ExplorerWeb.Batches.Index ), do: "text-foreground/80 hover:text-foreground font-semibold"

  defp active_view_class(ExplorerWeb.Operators.Index, ExplorerWeb.Operators.Index ), do: "text-green-500 font-bold"
  defp active_view_class(ExplorerWeb.Operator.Index, ExplorerWeb.Operators.Index ), do: "text-green-500 font-bold"
  defp active_view_class(_other, ExplorerWeb.Operators.Index ), do: "text-foreground/80 hover:text-foreground font-semibold"

end
