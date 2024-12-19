import { alignedTooltip } from "./tooltip";

export const costPerProofCustomOptions = (options, data) => {
	// add USD suffix
	options.scales.y.ticks.callback = (_value, index, values) => {
		if (index === 0) return `${Math.min(...data.datasets[0].data)} USD`;
		if (index === values.length - 1) {
			return `${Math.max(...data.datasets[0].data)} USD`;
		}
		return "";
	};

	options.plugins.tooltip.external = (context) =>
		alignedTooltip(context, {
			title: "Cost per proof",
			items: [
				{ title: "Cost", id: "cost" },
				{ title: "Age", id: "age" },
				{ title: "Merkle root", id: "merkle_root" },
				{ title: "Block number", id: "block_number" },
				{ title: "Amount of proofs", id: "amount_of_proofs" },
			],
			onTooltipUpdate: (tooltipModel) => {
				const dataset = tooltipModel.dataPoints[0].dataset;
				const idx = tooltipModel.dataPoints[0].dataIndex;

				const cost = `${dataset.data[idx]} USD`;
				const age = tooltipModel.dataPoints[0].label;
				const merkleRootHash = dataset.merkle_root[idx];
				const merkle_root = `${merkleRootHash.slice(0, 6)}...${merkleRootHash.slice(
					merkleRootHash.length - 4
				)}`;
				const block_number = dataset.submission_block_number[idx];
				const amount_of_proofs = dataset.amount_of_proofs[idx];

				return {
					cost,
					age,
					merkle_root,
					block_number,
					amount_of_proofs,
				};
			},
		});
};
