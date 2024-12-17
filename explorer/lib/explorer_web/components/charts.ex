defmodule ExplorerWeb.ChartComponents do
  use Phoenix.Component

  attr(:id, :string, required: true)
  attr(:chart_type, :string, required: true)
  attr(:chart_data, :string, required: true)
  attr(:chart_options, :string, required: true)

  defp basic_chart(assigns) do
    ~H"""
    <div class="relative h-full w-full">
      <canvas
        id={@id}
        phx-hook="ChartHook"
        data-chart-type={@chart_type}
        data-chart-data={@chart_data}
        data-chart-options={@chart_options}
        style="height: 100%; width: 100%;"
      >
      </canvas>
    </div>
    """
  end

  @doc """
  Renders a line chart with aligned style.

  ## Examples
      <.line id="exchanges" data=[1, 2, 3, 4] labels=["January", "February", "March", "April"]  ./>
  """
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
          maintainAspectRatio: false,
          interaction: %{
            mode: "nearest",
            intersect: false
          },
          plugins: %{
            legend: %{
              display: false
            }
          },
          scales: %{
            x: %{
              ticks: %{
                display: false
              },
              grid: %{
                display: false
              },
              border: %{
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
