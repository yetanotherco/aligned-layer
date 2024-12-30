export const findClosestIndex = (target, values) => {
	let closestIndex = 0;
	let smallestDiff = Math.abs(values[0] - target);
	for (let i = 1; i < values.length; i++) {
		const diff = Math.abs(values[i] - target);
		if (diff < smallestDiff) {
			closestIndex = i;
			smallestDiff = diff;
		}
	}
	return closestIndex;
};

export const yTickCallbackShowMinAndMaxValues =
	(data, renderText) => (_value, index, values) => {
		if (index === 0) return `0 USD`;

		const dataY = data.datasets[0].data.map((point) => parseFloat(point.y));
		const sortedData = dataY.sort((a, b) => b - a);
		const min = sortedData[0];
		const max = sortedData[sortedData.length - 1];
		const valsData = values.map((item) => item.value);
		const idxClosestToMin = findClosestIndex(min, valsData);
		const idxClosestToMax = findClosestIndex(max, valsData);

		if (index == idxClosestToMin) return renderText(min);
		if (index == idxClosestToMax) return renderText(max);
		return "";
	};
