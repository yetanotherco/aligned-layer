defmodule LabeledProgressBarComponent do
  use ExplorerWeb, :live_component

  attr :progress, :float, required: true
  attr :label, :string, required: true

  @impl true
  def render(assigns) do
  ~H"""
  <div class="progress-with-labels">
    <div
    class="progress-foreground"
    style={"clip-path: inset(0 #{100 - @progress}% 0 0);"}
    >
    <%= @label %>
    </div>
    <div
    class="progress-background"
    style={"clip-path: inset(0 0 0 #{@progress}%);"}
    >
    <%= @label %>
    </div>
  </div>
  """
  end

  def update_assigns(assigns, socket) do
    progress = assigns[:progress] || 0
    label = assigns[:label] || ""
    {:ok, assign(socket, progress: progress, label: label)}
  end
end
