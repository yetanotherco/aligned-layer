const tooltipItem = (title, valueId) => `
  <div class="chart-tooltip-item">
    <p class="chart-tooltip-item-title">${title}</p>
    <p class="chart-tooltip-item-value" id=${valueId}></p>
  </div>
`;

const tooltipComponent = ({ title, items }) => `
  <div class="chart-tooltip-container">
    <p class="chart-tooltip-title">${title}</p>
    <div class="chart-tooltip-items">
      ${items.map((item) => tooltipItem(item.title, item.id)).join("")}
    </div>
  </div>
`;

/**
 * Displays and positions a custom tooltip on the page based on the provided chart context.
 * This function creates the tooltip element if it doesn't exist, updates its content,
 * and positions it correctly based on the chart's canvas and tooltip model.
 *
 * @param {Object} context - The chart.js context, typically passed as `this` within chart hooks.
 * @param {Object} params - An object containing configuration for the tooltip.
 * @param {string} params.title - The title text to display in the tooltip.
 * @param {Array} params.items - An array of items (with ids) to be displayed inside the tooltip.
 * @param {Function} params.onTooltipUpdate - A callback function that updates the tooltip values based on the tooltip model.
 *                                           It is called with the `tooltipModel` as an argument and should return an object containing updated values for each tooltip item by their respective `id`.
 *
 * @returns {void} - This function doesn't return any value. It manipulates the DOM directly to show the tooltip.
 *
 * @example
 * alignedTooltip(context, {
 *   title: "Tooltip Title",
 *   items: [{ title: "Cost", id: "cost_id" }, { title: "Timestamp", id: "timestamp_id" }],
 *   onTooltipUpdate: (tooltipModel) => {
 * 	   const cost = tooltipModel.dataPoints[0].raw;
 * 	   const timestamp = tooltipModel.dataPoints[0].label;
 *
 *     return {
 *       cost_id: cost,
 *       timestamp_id: label
 *     };
 *   }
 * });
 */
export const alignedTooltip = (context, { title, items, onTooltipUpdate }) => {
	let tooltipEl = document.getElementById("chartjs-tooltip");
	if (!tooltipEl) {
		tooltipEl = document.createElement("div");
		tooltipEl.style = "transition: opacity 0.3s;";
		tooltipEl.style = "transition: left 0.1s;";
		tooltipEl.id = "chartjs-tooltip";
		tooltipEl.innerHTML = tooltipComponent({ title, items });
		document.body.appendChild(tooltipEl);
	}

	const tooltipModel = context.tooltip;
	// Hide element if no tooltip
	if (tooltipModel.opacity === 0) {
		tooltipEl.style.opacity = 0;
		return;
	}
	tooltipEl.classList.remove("above", "below", "no-transform");
	if (tooltipModel.yAlign) {
		tooltipEl.classList.add(tooltipModel.yAlign);
	} else {
		tooltipEl.classList.add("no-transform");
	}

	const values = onTooltipUpdate(tooltipModel);
	items.forEach((item) => {
		tooltipEl.querySelector(`#${item.id}`).textContent = values[item.id];
	});

	const position = context.chart.canvas.getBoundingClientRect();

	tooltipEl.style.opacity = 1;
	tooltipEl.style.position = "absolute";
	tooltipEl.style.left =
		position.left -
		tooltipEl.offsetWidth / 2 +
		window.scrollX +
		tooltipModel.caretX +
		"px";
	tooltipEl.style.top =
		position.top + window.scrollY + tooltipModel.caretY + "px";
	tooltipEl.style.pointerEvents = "none";
};
