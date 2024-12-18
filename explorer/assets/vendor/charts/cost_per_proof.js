export const costPerProofCustomOptions = (options, data) => {
	// y axis
	options.scales.y.ticks.display = true;
	options.scales.y.ticks.callback = (value, index, values) => {
		if (index === 0) return `${Math.min(...data.datasets[0].data)} USD`;
		if (index === values.length - 1)
			return `${Math.max(...data.datasets[0].data)} USD`;
		return "";
	};

	// x axis
	options.scales.x.ticks.display = true;
	options.scales.x.ticks.callback = (_, index) => {
		const formatDate = (isoString) => {
			const date = new Date(isoString);
			return date.toLocaleDateString("en-US", {
				month: "short",
				day: "numeric",
			});
		};

		if (index === 0) return formatDate(data.labels[0]);
		if (index === Math.floor((data.labels.length - 1) / 2))
			return formatDate(
				data.labels[Math.floor((data.labels.length - 1) / 2)]
			);
		if (index === data.labels.length - 1)
			return formatDate(data.labels[data.labels.length - 1]);

		return "";
	};
};
