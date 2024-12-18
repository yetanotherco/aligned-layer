defmodule ExplorerWeb.ChartComponents do
  use Phoenix.Component

  attr(:id, :string, required: true)
  attr(:chart_type, :string, required: true)
  attr(:chart_data, :string, required: true)
  attr(:chart_options, :string, required: true)
  attr(:chart_tooltip, :string, required: true)

  defp basic_chart(assigns) do
    ~H"""
    <div class="relative h-full w-full">
      <canvas
        id={@id}
        phx-hook="ChartHook"
        data-chart-type={@chart_type}
        data-chart-data={@chart_data}
        data-chart-options={@chart_options}
        data-chart-tooltip={@chart_tooltip}
        style="height: 100%; width: 100%;"
      >
      </canvas>
    </div>
    """
  end

  @doc """
  Renders a line chart with aligned style.

  ## Examples
    <.line
      id="exchanges"
      data={[1, 2, 3, 4]}
      labels={["January", "February", "March", "April"]}
      tooltip={%{title: "Exchange details", body: "Month: {{label}}\nRate: {{value}}"}}
    />
    !Note that {{label}} and {{value}} will get replaced with their respective values, the alternative would be to pass raw JS...
  """
  attr(:id, :string, required: true)
  attr(:labels, :list, required: true)
  attr(:data, :list, required: true)
  attr(:tooltip, :map, required: true)

  def line_chart(assigns) do
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
      chart_tooltip={Jason.encode!(@tooltip)}
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
          elements: %{
            point: %{
              pointStyle: false
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
