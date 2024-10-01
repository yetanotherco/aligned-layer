defmodule TelemetryApiWeb.TraceController do
  use TelemetryApiWeb, :controller

  alias TelemetryApi.Traces

  action_fallback TelemetryApiWeb.FallbackController

  @doc """
  Create a trace for a NewTask with the given merkle_root
  Method: POST initTaskTrace
  """
  def create_task_trace(conn, %{"merkle_root" => merkle_root}) do
    with :ok <- Traces.create_task_trace(merkle_root) do
      conn
      |> put_status(:created)
      |> render(:show_merkle, merkle_root: merkle_root)
    end
  end

  @doc """
  Register an operator response in the trace of the given merkle_root
  Method: POST operatorResponse
  """
  def register_operator_response(conn, %{
        "merkle_root" => merkle_root,
        "operator_id" => operator_id
      }) do
    with :ok <- Traces.register_operator_response(merkle_root, operator_id) do
      conn
      |> put_status(:ok)
      |> render(:show_operator, operator_id: operator_id)
    end
  end

  @doc """
  Registers a reached quorum in the trace of the given merkle_root
  Method: POST quorumReached
  """
  def quorum_reached(conn, %{"merkle_root" => merkle_root}) do
    with :ok <- Traces.quorum_reached(merkle_root) do
      conn
      |> put_status(:ok)
      |> render(:show_merkle, merkle_root: merkle_root)
    end
  end

  @doc """
  Registers an error in the trace of the given merkle_root
  Method: POST taskError
  """
  def task_error(conn, %{"merkle_root" => merkle_root, "error" => error}) do
    with :ok <- Traces.task_error(merkle_root, error) do
      conn
      |> put_status(:ok)
      |> render(:show_merkle, merkle_root: merkle_root)
    end
  end

  @doc """
  Finish a trace for the given merkle_root
  Method: POST finishTaskTrace
  """
  def finish_task_trace(conn, %{"merkle_root" => merkle_root}) do
    with :ok <- Traces.finish_task_trace(merkle_root) do
      conn
      |> put_status(:ok)
      |> render(:show_merkle, merkle_root: merkle_root)
    end
  end
end
