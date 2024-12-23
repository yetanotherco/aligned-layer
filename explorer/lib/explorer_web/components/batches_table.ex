defmodule ExplorerWeb.BatchesTable do
  use Phoenix.Component
  use ExplorerWeb, :live_component

  attr(:batches, :list, required: true)

  def batches_table(assigns) do
    ~H"""
    <.table id="batches" rows={@batches}>
      <:col :let={latest_batch} label="Batch Hash" class="text-left">
        <.link navigate={~p"/batches/#{latest_batch.merkle_root}"}>
          <span class="inline-flex gap-x-3 items-center group-hover:text-foreground/80">
            <%= Helpers.shorten_hash(latest_batch.merkle_root, 6) %>
            <.right_arrow />
            <.tooltip>
              <%= latest_batch.merkle_root %>
            </.tooltip>
          </span>
        </.link>
      </:col>
      <:col :let={latest_batch} label="Status">
        <.dynamic_badge_for_batcher status={Helpers.get_batch_status(latest_batch)} />
      </:col>
      <:col :let={latest_batch} label="Age">
        <span class="md:px-0" title={latest_batch.submission_timestamp}>
          <%= latest_batch.submission_timestamp |> Helpers.parse_timeago() %>
        </span>
      </:col>
      <:col :let={latest_batch} label="Block Number">
        <%= latest_batch.submission_block_number |> Helpers.format_number() %>
      </:col>
    </.table>
    """
  end
end
