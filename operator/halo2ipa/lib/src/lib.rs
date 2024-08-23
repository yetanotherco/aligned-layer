use halo2_proofs::{
    plonk::{read_fr, read_params, verify_proof, VerifyingKey},
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
use std::{io::BufReader, slice};

#[no_mangle]
pub extern "C" fn verify_halo2_ipa_proof_ffi(
    proof_buf: *const u8,
    proof_len: u32,
    params_buf: *const u8,
    params_len: u32,
    public_input_buf: *const u8,
    public_input_len: u32,
) -> bool {
    if proof_buf.is_null() || params_buf.is_null() || public_input_buf.is_null() {
        error!("Halo2 IPA input buffer length null");
        return false;
    }

    // For Halo2 the `params_buf` contains the serialized cs, vk, and params with there respective sizes serialized as u32 values (4 bytes) => 3 * 4 bytes = 12 by the concatenated variable length buffers:
    // We therefore require that the `params_buf` is greater than 12 bytes and treat the case that buffer lengths and buffers themselves are 0 size as false.
    // [ cs_len | vk_len | vk_params_len | cs_bytes | vk_bytes | vk_params_bytes ].
    if proof_len == 0 || params_len <= 12 || public_input_len == 0 {
        error!("Halo2 IPA input buffer length zero size");
        return false;
    }

    let proof_bytes = unsafe { slice::from_raw_parts(proof_buf, proof_len as usize) };

    let params_bytes = unsafe { slice::from_raw_parts(params_buf, params_len as usize) };

    let public_input_bytes =
        unsafe { slice::from_raw_parts(public_input_buf, public_input_len as usize) };

    // Deserialize bytes
    let Ok((cs_bytes, vk_bytes, vk_params_bytes)) = read_params(params_bytes) else {
        error!("Failed to deserialize verifiation parameter buffers from parameters buffer");
        return false;
    };

    let Ok(cs) = bincode::deserialize(cs_bytes) else {
        error!("Failed to deserialize constraint system");
        return false;
    };

    let Ok(vk) =
        VerifyingKey::<G1Affine>::read(&mut BufReader::new(vk_bytes), SerdeFormat::RawBytes, cs)
    else {
        error!("Failed to deserialize verification key");
        return false;
    };

    let Ok(params) = Params::read::<_>(&mut BufReader::new(vk_params_bytes)) else {
        error!("Failed to deserialize verification parameters");
        return false;
    };

    let Ok(res) = read_fr(public_input_bytes) else {
        error!("Failed to deserialize public inputs");
        return false;
    };
    let strategy = SingleStrategy::new(&params);
    let instances = res;
    let mut transcript = Blake2bRead::<&[u8], G1Affine, Challenge255<_>>::init(proof_bytes);
    verify_proof::<
        IPACommitmentScheme<G1Affine>,
        VerifierIPA<G1Affine>,
        Challenge255<G1Affine>,
        Blake2bRead<&[u8], G1Affine, Challenge255<G1Affine>>,
        SingleStrategy<G1Affine>,
    >(&params, &vk, strategy, &[vec![instances]], &mut transcript)
    .is_ok()
}

#[cfg(test)]
mod tests {
    use super::*;

    use ff::{Field, PrimeField};
    use halo2_backend::poly::commitment::ParamsProver;
    use halo2_proofs::{
        plonk::{
            create_proof, keygen_pk, keygen_vk_custom, verify_proof, write_params, StandardPlonk,
        },
        poly::ipa::{
            commitment::{IPACommitmentScheme, ParamsIPA},
            multiopen::ProverIPA,
        },
        transcript::{
            Blake2bRead, Blake2bWrite, Challenge255, TranscriptReadBuffer, TranscriptWriterBuffer,
        },
    };
    use halo2curves::bn256::Fr;
    use rand_core::OsRng;
    use std::{
        fs::File,
        io::{BufWriter, Read, Write},
    };

    const PROOF: &[u8] = include_bytes!("../../../../scripts/test_files/halo2_ipa/proof.bin");

    const PUB_INPUT: &[u8] =
        include_bytes!("../../../../scripts/test_files/halo2_ipa/pub_input.bin");

    const PARAMS: &[u8] = include_bytes!("../../../../scripts/test_files/halo2_ipa/params.bin");

    #[test]
    fn halo2_serialization_works() {
        let k = 4;
        let circuit = StandardPlonk(Fr::random(OsRng));
        let params = ParamsIPA::<G1Affine>::new(k);
        let compress_selectors = true;
        let vk =
            keygen_vk_custom(&params, &circuit, compress_selectors).expect("vk should not fail");
        let cs = vk.clone().cs;
        let instances = vec![vec![circuit.0]];
        let pk = keygen_pk(&params, vk.clone(), &circuit).expect("pk should not fail");

        let mut transcript = Blake2bWrite::<_, _, Challenge255<_>>::init(vec![]);
        create_proof::<
            IPACommitmentScheme<G1Affine>,
            ProverIPA<G1Affine>,
            Challenge255<G1Affine>,
            _,
            Blake2bWrite<Vec<u8>, G1Affine, Challenge255<_>>,
            _,
        >(
            &params,
            &pk,
            &[circuit.clone()],
            &[instances.clone()],
            OsRng,
            &mut transcript,
        )
        .expect("prover should not fail");

        let proof = transcript.finalize();
        let strategy = SingleStrategy::new(&params);
        let mut transcript = Blake2bRead::<_, _, Challenge255<_>>::init(&proof[..]);

        verify_proof::<
            IPACommitmentScheme<G1Affine>,
            VerifierIPA<G1Affine>,
            Challenge255<G1Affine>,
            Blake2bRead<&[u8], G1Affine, Challenge255<G1Affine>>,
            SingleStrategy<G1Affine>,
        >(
            &params,
            &vk,
            strategy,
            &[instances.clone()],
            &mut transcript,
        )
        .expect("verifier shoud not fail");

        //write proof
        std::fs::write("proof.bin", &proof[..]).expect("should succeed to write new proof");

        //write public input
        let f = File::create("pub_input.bin").unwrap();
        let mut writer = BufWriter::new(f);
        instances.to_vec().into_iter().flatten().for_each(|fp| {
            writer.write(&fp.to_repr()).unwrap();
        });
        writer.flush().unwrap();

        let mut vk_buf = Vec::new();
        vk.write(&mut vk_buf, SerdeFormat::RawBytes).unwrap();

        let mut params_buf = Vec::new();
        params.write(&mut params_buf).unwrap();

        //Write everything to parameters file
        File::create("params.bin").unwrap();
        write_params(&params_buf, cs, &vk_buf, "params.bin").unwrap();

        //read proof
        let proof = std::fs::read("proof.bin").expect("should succeed to read proof");

        //read instances
        let mut f = File::open("pub_input.bin").unwrap();
        let mut buf = Vec::new();
        f.read_to_end(&mut buf).unwrap();
        let res = read_fr(&buf).unwrap();
        let instances = res;

        let mut f = File::open("params.bin").unwrap();
        let mut params_buf = Vec::new();
        f.read_to_end(&mut params_buf).unwrap();

        let (cs_bytes, vk_bytes, vk_params_bytes) = read_params(&params_buf).unwrap();

        let cs = bincode::deserialize(cs_bytes).unwrap();
        let vk = VerifyingKey::<G1Affine>::read(
            &mut BufReader::new(vk_bytes),
            SerdeFormat::RawBytes,
            cs,
        )
        .unwrap();
        let params = Params::read::<_>(&mut BufReader::new(vk_params_bytes)).unwrap();

        let strategy = SingleStrategy::new(&params);
        let mut transcript = Blake2bRead::<_, _, Challenge255<_>>::init(&proof[..]);
        assert!(verify_proof::<
            IPACommitmentScheme<G1Affine>,
            VerifierIPA<G1Affine>,
            Challenge255<G1Affine>,
            Blake2bRead<&[u8], G1Affine, Challenge255<G1Affine>>,
            SingleStrategy<G1Affine>,
        >(&params, &vk, strategy, &[vec![instances]], &mut transcript)
        .is_ok());
        std::fs::remove_file("proof.bin").unwrap();
        std::fs::remove_file("pub_input.bin").unwrap();
        std::fs::remove_file("params.bin").unwrap();
    }

    #[test]
    fn verify_halo2_plonk_proof() {
        let proof_len = PROOF.len();
        let proof_bytes = PROOF.as_ptr();

        let params_len = PARAMS.len();
        let params_bytes = PARAMS.as_ptr();

        let public_input_len = PUB_INPUT.len();
        let public_input_bytes = PUB_INPUT.as_ptr();

        let result = verify_halo2_ipa_proof_ffi(
            proof_bytes,
            proof_len as u32,
            params_bytes,
            params_len as u32,
            public_input_bytes,
            public_input_len as u32,
        );
        assert!(result)
    }

    #[test]
    fn verify_halo2_plonk_proof_aborts_with_bad_proof() {
        // Select Proof Bytes
        let proof_len = PROOF.len();
        let proof_bytes = PROOF.as_ptr();

        let params_len = PARAMS.len();
        let params_bytes = PARAMS.as_ptr();

        let public_input_len = PUB_INPUT.len();
        let public_input_bytes = PUB_INPUT.as_ptr();

        let result = verify_halo2_ipa_proof_ffi(
            proof_bytes,
            (proof_len - 1) as u32,
            params_bytes,
            params_len as u32,
            public_input_bytes,
            public_input_len as u32,
        );
        assert!(!result)
    }
}
