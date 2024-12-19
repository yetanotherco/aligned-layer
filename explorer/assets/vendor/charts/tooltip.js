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
