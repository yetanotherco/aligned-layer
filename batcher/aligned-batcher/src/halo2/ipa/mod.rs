use halo2_proofs::{
    plonk::{verify_proof, VerifyingKey, read_params, read_fr},
    poly::{
        commitment::Params,
        ipa::{commitment::IPACommitmentScheme, multiopen::VerifierIPA, strategy::SingleStrategy},
        VerificationStrategy,
    },
    transcript::{Blake2bRead, Challenge255, TranscriptReadBuffer},
    SerdeFormat,
};
use halo2curves::bn256::G1Affine;
use log::error;
use std::io::BufReader;

// MaxConstraintSystemSize 2KB
pub const MAX_CONSTRAINT_SYSTEM_SIZE: usize = 2 * 1024;

// MaxVerificationKeySize 1KB
pub const MAX_VERIFIER_KEY_SIZE: usize = 1024;

// MaxipaParamsSize 4KB
pub const MAX_IPA_PARAMS_SIZE: usize = 4 * 1024;

pub fn verify_halo2_ipa(proof: &[u8], public_input: &[u8], verification_key: &[u8]) -> bool {
    // For Halo2 the `verification_key` contains the serialized cs, vk, and params with there respective sizes serialized as u32 values (4 bytes) => 3 * 4 bytes = 12:
    // We therefore require that the `verification_key` is greater than 12 bytes and treat the case that buffer lengths and buffers themselves are 0 size as false.
    // [ cs_len | vk_len | vk_params_len | cs_bytes | vk_bytes | vk_params_bytes ].
    if verification_key.len() <= 12 {
        error!("verification input buffers less than 12 bytes");
        return false;
    } else if proof.is_empty() {
        error!("proof input buffers zero size");
        return false;
    } else if public_input.is_empty() {
        error!("public input input buffers zero size");
        return false;
    }

    let Ok((cs_bytes, vk_bytes, vk_params_bytes)) = read_params(verification_key) else {
        error!("Failed to deserialize verifiation parameter buffers from parameters buffer");
        return false;
    };

    let Ok(cs) = bincode::deserialize(&cs_bytes) else {
        error!("Failed to deserialize constraint system");
        return false;
    };
    
    let Ok(vk) = VerifyingKey::<G1Affine>::read(
            &mut BufReader::new(vk_bytes),
            SerdeFormat::RawBytes,
            cs,
        ) else {
        error!("Failed to deserialize verification key");
        return false;
    };

    let Ok(params) = Params::read::<_>(&mut BufReader::new(vk_params_bytes)) else {
        error!("Failed to deserialize verification parameters");
        return false;
    };
    
    let Ok(res) = read_fr(public_input) else {
        error!("Failed to deserialize public inputs");
        return false;
    };

    let strategy = SingleStrategy::new(&params);
    let instances = res;
    let mut transcript =
        Blake2bRead::<&[u8], G1Affine, Challenge255<_>>::init(proof);
    verify_proof::<
        IPACommitmentScheme<G1Affine>,
        VerifierIPA<G1Affine>,
        Challenge255<G1Affine>,
        Blake2bRead<&[u8], G1Affine, Challenge255<G1Affine>>,
        SingleStrategy<G1Affine>,
    >(
        &params, &vk, strategy, &[vec![instances]], &mut transcript
    )
    .is_ok()
}