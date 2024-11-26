defmodule BLSSignatureVerifier do
  def verify(signature, {pubkey_g1_x, pubkey_g1_y}, bls_pubkey_g2, message) do
    endian = :big
    pubkey_g1_x = <<pubkey_g1_x::unsigned-big-integer-size(256)>>
    pubkey_g1_y = <<pubkey_g1_y::unsigned-big-integer-size(256)>>

    args = [
      "--signature",
      :binary.list_to_bin(signature),
      "--public-key-g1-x",
      Base.encode16(pubkey_g1_x),
      "--public-key-g1-y",
      Base.encode16(pubkey_g1_y),
      "--public-key-g2",
      :binary.list_to_bin(bls_pubkey_g2),
      "--message",
      Base.encode16(message)
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
