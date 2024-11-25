defmodule BLSSignatureVerifier do
  def verify(signature, public_key, message) do
    # Encode the args as hex
    args = [
      "--signature",
      Base.encode16(signature, case: :lower),
      "--publickey",
      Base.encode16(public_key, case: :lower),
      "--message",
      Base.encode16(message, case: :lower)
    ]

    binary_path = Path.join(:code.priv_dir(:telemetry_api), "bls_verify")
    {output, exit_code} = System.cmd(binary_path, args)

    case exit_code do
      0 -> {:ok, "Valid"}
      1 -> {:error, "Invalid signature"}
      _ -> {:error, "Verification failed: #{output}"}
    end
  end
end
