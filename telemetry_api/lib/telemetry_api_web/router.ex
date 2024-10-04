defmodule TelemetryApiWeb.Router do
  alias TelemetryApiWeb.OperatorController
  use TelemetryApiWeb, :router

  pipeline :api do
    plug :accepts, ["json"]
  end

  scope "/api", TelemetryApiWeb do
    pipe_through :api
    post "/operators", OperatorController, :create_or_update
    resources "/operators", OperatorController, only: [:index, :show]
  end

  scope "/versions", TelemetryApiWeb do
    pipe_through :api
    post "/operators", OperatorController, :create_or_update
    resources "/", OperatorController, only: [:index, :show]
  end

  # Enable LiveDashboard in development
  # if Application.compile_env(:telemetry_api, :dev_routes) do
  #   # If you want to use the LiveDashboard in production, you should put
  #   # it behind authentication and allow only admins to access it.
  #   # If your application does not have an admins-only section yet,
  #   # you can use Plug.BasicAuth to set up some basic authentication
  #   # as long as you are also using SSL (which you should anyway).
  #   import Phoenix.LiveDashboard.Router

  #   scope "/dev" do
  #     pipe_through [:fetch_session, :protect_from_forgery]

  #     live_dashboard "/dashboard", metrics: TelemetryApiWeb.Telemetry
  #   end
  # end
end
