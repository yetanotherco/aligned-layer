defmodule BLSApkRegistry do
  require Logger

  @aligned_config_file System.get_env("ALIGNED_CONFIG_FILE")

  config_file_path =
    case @aligned_config_file do
      nil -> raise("ALIGNED_CONFIG_FILE not set in .env")
      file -> file
    end

  {status, config_json_string} = File.read(config_file_path)

  case status do
    :ok ->
      Logger.debug("Aligned deployment file read successfully")

    :error ->
      raise(
        "Config file not read successfully, did you run make explorer_create_env? If you did,\n make sure Eigenlayer config file is correctly stored"
      )
  end

  @contract_address Jason.decode!(config_json_string)
                    |> Map.get("addresses")
                    |> Map.get("blsApkRegistry")

  use Ethers.Contract,
    abi_file: "priv/abi/IBLSApkRegistry.json",
    default_address: @contract_address

  def get_registry_coordinator_address() do
    @contract_address
  end

  def get_operator_bls_pubkey(operator_address) do
    case BLSApkRegistry.get_registered_pubkey(operator_address)
         |> Ethers.call() do
      {:ok, data} ->
        data

      error ->
        {:error, error}
    end
  end
end
