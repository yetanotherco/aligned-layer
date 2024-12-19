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
			],
			onTooltipUpdate: (tooltipModel) => {
				const value = tooltipModel.dataPoints[0].raw;
				const label = tooltipModel.dataPoints[0].label;

				return {
					cost: value,
					age: label,
				};
			},
		});
};
