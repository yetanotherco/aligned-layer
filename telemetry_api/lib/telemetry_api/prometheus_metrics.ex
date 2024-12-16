defmodule TelemetryApi.PrometheusMetrics do
  use Prometheus.Metric

  @gauge [name: :gas_price, help: "Ethereum Gas Price.", labels: []]
  @counter [name: :missing_operator_count, help: "Missing Operators", labels: [:operator]]

  def new_gas_price(gas_price) do
    Gauge.set(
      [name: :gas_price, labels: []],
      gas_price
    )
  end

  def missing_operator(operator) do
    Counter.inc(
      name: :missing_operator_count,
      labels: [operator]
    )
  end
end
