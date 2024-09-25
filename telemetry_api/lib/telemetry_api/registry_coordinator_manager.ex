defmodule TelemetryApi.RegistryCoordinatorManager do
  alias TelemetryApi.RegistryCoordinatorManager

  @registry_coordinator_address System.get_env("REGISTRY_COORDINATOR_ADDRESS")

  use Ethers.Contract,
    abi_file: "lib/abi/IRegistryCoordinator.json",
    default_address: System.get_env("REGISTRY_COORDINATOR_ADDRESS")

  def get_registry_coordinator_address() do
    @registry_coordinator_address
  end

  def is_operator_registered?(operator_address) do
    {:ok, operator_address } = Base.decode16(operator_address, case: :mixed)

    case RegistryCoordinatorManager.get_operator_status(operator_address)
      |> Ethers.call() do
        {:ok, response} -> response == 1
        error ->
          {:error, error}
      end
  end
end
