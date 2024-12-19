import { alignedTooltip } from "./tooltip";

const getDateInTimesAgo = (isoDate) => {
	const now = new Date();
	const date = new Date(isoDate);
	const diffInSeconds = Math.floor((now - date) / 1000);

	const intervals = [
		{ label: "year", seconds: 31536000 },
		{ label: "month", seconds: 2592000 },
		{ label: "week", seconds: 604800 },
		{ label: "day", seconds: 86400 },
		{ label: "hour", seconds: 3600 },
		{ label: "minute", seconds: 60 },
		{ label: "second", seconds: 1 },
	];

	for (const interval of intervals) {
		const count = Math.floor(diffInSeconds / interval.seconds);
		if (count >= 1) {
			return `${count} ${interval.label}${count > 1 ? "s" : ""} ago`;
		}
	}

	return "just now";
};

export const costPerProofCustomOptions = (options, data) => {
	// y axis
	options.scales.y.ticks.display = true;
	options.scales.y.ticks.callback = (_value, index, values) => {
		if (index === 0) return `${Math.min(...data.datasets[0].data)} USD`;
		if (index === values.length - 1) {
			return `${Math.max(...data.datasets[0].data)} USD`;
		}
		return "";
	};

	// x axis
	options.scales.x.ticks.display = true;
	options.scales.x.ticks.callback = (_value, index, values) => {
		if (index === 0) return data.labels[0];
		if (index === Math.floor((data.labels.length - 1) / 2))
			return data.labels[Math.floor((data.labels.length - 1) / 2)];
		if (index === values.length - 1)
			return data.labels[data.labels.length - 1];

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
