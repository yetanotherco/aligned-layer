defmodule ExplorerWeb.ChartComponents do
  use Phoenix.Component

  @doc """
  Renders a line chart with aligned style.

  ## Examples
      <.ChartComponents.line datasets="" ./>
  """
  attr(:id, :string, required: true)
  attr(:chart_type, :string, required: true)
  attr(:chart_data, :string, required: true)
  attr(:chart_options, :string, required: true)

  def basic_chart(assigns) do
    ~H"""
    <div class="chart-container" class="relative w-full h-full">
      <canvas
        id={@id}
        phx-hook="ChartHook"
        data-chart-type={@chart_type}
        data-chart-data={@chart_data}
        data-chart-options={@chart_options}
        class="w-full"
      >
      </canvas>
    </div>
    """
  end

  attr(:id, :string, required: true)
  attr(:labels, :list, required: true)
  attr(:data, :list, required: true)

  def line(assigns) do
    ~H"""
    <.basic_chart
      id={@id}
      chart_type="line"
      chart_data={
        Jason.encode!(%{
          labels: @labels,
          datasets: [%{data: @data, borderColor: "rgb(24, 255, 127)", fill: false, tension: 0.1}]
        })
      }
      chart_options={
        Jason.encode!(%{
          plugins: %{
            legend: %{
              display: false
            }
          },
          scales: %{
            x: %{
              grid: %{
                display: false
              }
            },
            y: %{
              ticks: %{
                display: false
              },
              grid: %{
                display: false
              },
              border: %{
                display: false
              }
            }
          }
        })
      }
    />
    """
  end
end
