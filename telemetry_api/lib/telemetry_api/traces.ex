defmodule TelemetryApi.Traces do
  @moduledoc """
  The Traces context.
  """
  alias TelemetryApi.Traces.Trace
  alias TelemetryApi.Operators
  alias TelemetryApi.ContractManagers.StakeRegistry

  require OpenTelemetry.Tracer
  require OpenTelemetry.Ctx
  alias OpenTelemetry.Tracer, as: Tracer
  alias OpenTelemetry.Ctx, as: Ctx

  @doc """
  Send the trace to OpenTelemetry

  This function is responsible for creating a new span and storing the context in the Agent.

  ## Examples

      iex> merkle_root = "0x1234567890abcdef"
      iex> create_task_trace(merkle_root)
      :ok
  """
  def create_task_trace(merkle_root) do
    with {:ok, trace} <- set_current_trace(merkle_root) do
      with {:ok, total_stake} <- StakeRegistry.get_current_total_stake() do
        aggregator_subspan_ctx =
          Tracer.start_span(
            "Aggregator",
            %{
              attributes: [
                {:merkle_root, merkle_root},
                {:total_stake, total_stake}
              ]
            }
          )

        Tracer.set_current_span(aggregator_subspan_ctx)
        Tracer.add_event("New task event received", [])

        TraceStore.store_trace(merkle_root, %{
          trace
          | subspans: Map.put(trace.subspans, :aggregator, aggregator_subspan_ctx)
        })

        :ok
      end
    end
  end

  @doc """
  Registers an operator response in the task trace.

  ## Examples

      iex> merkle_root = "0x1234567890abcdef"
      iex> operator_id = "0x..."
      iex> register_operator_response(merkle_root, operator_id)
      :ok
  """
  def register_operator_response(merkle_root, operator_id) do
    with {:ok, operator} <- Operators.get_operator(%{id: operator_id}),
         :ok <- validate_operator_registration(operator),
         {:ok, trace} <- set_current_trace_with_subspan(merkle_root, :aggregator) do
      operator_stake = Decimal.new(operator.stake)
      new_stake = Decimal.add(trace.current_stake, operator_stake)
      new_stake_fraction = Decimal.div(new_stake, trace.total_stake)
      operator_stake_fraction = Decimal.div(operator_stake, trace.total_stake)

      Tracer.add_event(
        "Operator Response: " <> operator.name,
        [
          {:merkle_root, merkle_root},
          {:operator_id, operator_id},
          {:name, operator.name},
          {:address, operator.address},
          {:operator_stake, Decimal.to_string(operator_stake)},
          {:current_stake, Decimal.to_string(new_stake)},
          {:current_stake_fraction, Decimal.to_string(new_stake_fraction)},
          {:operator_stake_fraction, Decimal.to_string(operator_stake_fraction)}
        ]
      )

      responses = trace.responses ++ [operator_id]

      TraceStore.store_trace(merkle_root, %{
        trace
        | responses: responses,
          current_stake: new_stake
      })

      IO.inspect(
        "Operator response included. merkle_root: #{IO.inspect(merkle_root)} operator_id: #{IO.inspect(operator_id)}"
      )

      :ok
    end
  end

  @doc """
  Registers the failure creating a batcher task in the task trace.

  ## Examples

      iex> merkle_root
      iex> error
      iex> batcher_task_creation_failed(merkle_root, error)
      :ok
  """
  def batcher_task_creation_failed(merkle_root, error) do
    with {:ok, trace} <- set_current_trace_with_subspan(merkle_root, :batcher) do
      Tracer.add_event(
        "Batcher Task Creation Failed",
        [
          {:error, error}
        ]
      )

      Tracer.end_span()

      TraceStore.store_trace(merkle_root, %{
        trace
        | subspans: Map.delete(trace.subspans, :batcher)
      })

      with {:ok, _trace} <- set_current_trace(merkle_root) do
        Tracer.end_span()
        TraceStore.delete_trace(merkle_root)
      end

      :ok
    end
  end

  @doc """
  Create a new task trace for the batcher and starts the subspan for the batcher.

  ## Examples

      iex> merkle_root
      iex> create_batcher_task_trace(merkle_root)
      :ok
  """
  def create_batcher_task_trace(merkle_root) do
    root_span_ctx =
      Tracer.start_span(
        "Task: #{merkle_root}",
        %{
          attributes: [
            {:merkle_root, merkle_root}
          ]
        }
      )

    {:ok, total_stake} = StakeRegistry.get_current_total_stake()
    ctx = Ctx.get_current()

    TraceStore.store_trace(merkle_root, %Trace{
      parent_span: root_span_ctx,
      context: ctx,
      total_stake: total_stake,
      current_stake: 0,
      responses: [],
      subspans: %{}
    })

    with {:ok, trace} <- set_current_trace(merkle_root) do
      # This span ends inmediately after it's created just to set the correct title to the final task.
      Tracer.with_span "Task: #{merkle_root}" do
        Tracer.set_attributes(%{merkle_root: merkle_root})
      end

      batcher_subspan_ctx =
        Tracer.start_span(
          "Batcher",
          %{
            attributes: [
              {:merkle_root, merkle_root}
            ]
          }
        )

      Tracer.set_current_span(batcher_subspan_ctx)
      Tracer.add_event("New batch", [{:merkle_root, merkle_root}])

      TraceStore.store_trace(merkle_root, %{
        trace
        | subspans: Map.put(trace.subspans, :batcher, batcher_subspan_ctx)
      })

      :ok
    end
  end

  @doc """
  Registers the uploading of a batcher task to S3 in the task trace.

  ## Examples

      iex> merkle_root
      iex> batcher_task_uploaded_to_s3(merkle_root)
      :ok
  """
  def batcher_task_uploaded_to_s3(merkle_root) do
    with {:ok, _trace} <- set_current_trace_with_subspan(merkle_root, :batcher) do
      Tracer.add_event("Batcher Task Uploaded to S3", [])
      :ok
    end
  end

  @doc """
  Registers the sending of a batcher task to Ethereum in the task trace.

  ## Examples

      iex> merkle_root
      iex> tx_hash
      iex> batcher_task_sent(merkle_root, tx_hash)
      :ok
  """
  def batcher_task_sent(merkle_root, tx_hash) do
    with {:ok, trace} <- set_current_trace_with_subspan(merkle_root, :batcher) do
      Tracer.add_event("Batcher Task Sent to Ethereum", [{"tx_hash", tx_hash}])

      Tracer.end_span()

      TraceStore.store_trace(merkle_root, %{
        trace
        | subspans: Map.delete(trace.subspans, :batcher)
      })

      :ok
    end
  end

  @doc """
  Registers the start of the creation of a batcher task in the task trace.

  ## Examples

      iex> merkle_root
      iex> batcher_task_started(merkle_root)
      :ok
  """
  def batcher_task_started(merkle_root, fee_per_proof, total_proofs) do
    with {:ok, _trace} <- set_current_trace_with_subspan(merkle_root, :batcher) do
      {fee_per_proof, _} = String.slice(fee_per_proof, 2..-1//1) |> Integer.parse(16)
      IO.inspect("fee_per_proof: #{fee_per_proof}")

      Tracer.add_event("Batcher Task being created",
        fee_per_proof: fee_per_proof,
        total_proofs: total_proofs
      )

      :ok
    end
  end

  @doc """
  Registers a reached quorum in the task trace.

  ## Examples

      iex> merkle_root = "0x1234567890abcdef"
      iex> quorum_reached(merkle_root)
      :ok
  """
  def quorum_reached(merkle_root) do
    with {:ok, _trace} <- set_current_trace_with_subspan(merkle_root, :aggregator) do
      Tracer.add_event("Quorum Reached", [])
      IO.inspect("Reached quorum registered. merkle_root: #{merkle_root}")
      :ok
    end
  end

  @doc """
  Registers an error in the task trace.

  ## Examples

      iex> merkle_root = "0x1234567890abcdef"
      iex> error = "Some error.."
      iex> task_error(merkle_root, error)
      :ok
  """
  def task_error(merkle_root, error) do
    with {:ok, _trace} <- set_current_trace(merkle_root) do
      Tracer.add_event(
        "Batch verification failed",
        [
          {:status, "error"},
          {:error, error}
        ]
      )

      IO.inspect("Task error registered. merkle_root: #{IO.inspect(merkle_root)}")
      :ok
    end
  end

  @doc """
  Registers the sending of an aggregator task to Ethereum in the task trace.

  ## Examples

      iex> merkle_root
      iex> tx_hash
      iex> aggregator_task_sent(merkle_root, tx_hash)
      :ok
  """
  def aggregator_task_sent(merkle_root, tx_hash) do
    with {:ok, _trace} <- set_current_trace_with_subspan(merkle_root, :aggregator) do
      Tracer.add_event("Task Sent to Ethereum", [{"tx_hash", tx_hash}])
      :ok
    end
  end

  @doc """
  Finish the task trace

  This function is responsible for ending the span and cleaning up the context.

  ## Examples

      iex> merkle_root = "0x1234567890abcdef"
      iex> finish_task_trace(merkle_root)
      :ok
  """
  def finish_task_trace(merkle_root) do
    with {:ok, trace} <- set_current_trace_with_subspan(merkle_root, :aggregator) do
      missing_operators =
        Operators.list_operators()
        |> Enum.filter(fn o -> o.id not in trace.responses and Operators.is_registered?(o) end)

      add_missing_operators(missing_operators)

      Tracer.end_span()

      with {:ok, _trace} <- set_current_trace(merkle_root) do
        Tracer.end_span()
        TraceStore.delete_trace(merkle_root)
      end

      # Clean up the context from the Agent
      IO.inspect("Finished task trace with merkle_root: #{merkle_root}.")
      :ok
    end
  end

  defp add_missing_operators([]), do: :ok

  defp add_missing_operators(missing_operators) do
    missing_operators =
      missing_operators |> Enum.map(fn o -> o.name end) |> Enum.join(";")

    Tracer.add_event("Missing Operators", [{:operators, missing_operators}])
  end

  defp set_current_trace(merkle_root) do
    with {:ok, trace} <- TraceStore.get_trace(merkle_root) do
      Ctx.attach(trace.context)
      Tracer.set_current_span(trace.parent_span)
      {:ok, trace}
    end
  end

  defp set_current_trace_with_subspan(merkle_root, span_name) do
    with {:ok, trace} <- TraceStore.get_trace(merkle_root) do
      Ctx.attach(trace.context)
      Tracer.set_current_span(trace.subspans[span_name])
      {:ok, trace}
    end
  end

  defp validate_operator_registration(operator) do
    if Operators.is_registered?(operator) do
      :ok
    else
      {:error, :bad_request, "Operator not registered"}
    end
  end
end
