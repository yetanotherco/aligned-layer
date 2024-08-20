use risc0_zkvm::{InnerReceipt, Receipt};
use log::error;

pub fn verify_risc_zero_proof(
    inner_receipt_bytes: &[u8],
    image_id: &[u8; 32],
    public_input: &[u8],
) -> bool {
    // We verify that the buffers are non-zero otherwise return false. We allow public_input size of 0.
    if receipt_bytes.is_empty() {
        error!("Risc0 receipt input buffer zero size");
        return false;
    } else if image_id.is_empty() {
        error!("Risc0 image_id input buffer zero size");
        return false;
    }

    if let Ok(inner_receipt) = bincode::deserialize::<InnerReceipt>(inner_receipt_bytes) {
        let receipt = Receipt::new(inner_receipt, public_input.to_vec());

        return receipt.verify(*image_id).is_ok();
    }
    false
}
