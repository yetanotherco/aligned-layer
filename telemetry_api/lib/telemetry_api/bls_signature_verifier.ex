defmodule BlsSignatureVerifier do
  use Rustler, otp_app: :bls_recover, crate: "blssignaturerecover"

  # Define the Rust NIF function
  def recover_public_key(signature, message_hash), do: :erlang.nif_error(:nif_not_loaded)
end
