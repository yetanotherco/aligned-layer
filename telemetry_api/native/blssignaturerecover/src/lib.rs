use bls12_381::{pairing, G1Projective, G2Projective, Scalar};
use rustler::{Encoder, Env, Term};
use sha2::{Digest, Sha256};

#[rustler::nif]
fn recover_public_key(signature: Vec<u8>, message_hash: Vec<u8>) -> Result<Vec<u8>, &'static str> {
    let signature = G1Projective::from_bytes(&signature).map_err(|_| "Invalid signature")?;
    let hash_scalar = Scalar::from_bytes(&message_hash).map_err(|_| "Invalid message hash")?;
    let public_key = G2Projective::generator() * hash_scalar;

    Ok(public_key.to_bytes().to_vec())
}

rustler::init!("Elixir.BlsSignatureVerifier");
