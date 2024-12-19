import Chart from "chart.js/auto";
import { costPerProofCustomOptions } from "./cost_per_proof";

const applyCommonChartOptions = (options, data) => {
	// generally, we want to show only the min and max values on the y axis
	if (options.scales.x.ticks.display) {
		options.scales.y.ticks.callback = (_value, index, values) => {
			if (index === 0) return Math.min(...data.datasets[0].data);
			if (index === values.length - 1) {
				return Math.max(...data.datasets[0].data);
			}
			return "";
		};
	}

	// on the x axes we want to show: the min, mean, max
	if (options.scales.y.ticks.display) {
		options.scales.x.ticks.callback = (_value, index, values) => {
			if (index === 0) return data.labels[0];
			if (index === Math.floor((data.labels.length - 1) / 2))
				return data.labels[Math.floor((data.labels.length - 1) / 2)];
			if (index === values.length - 1)
				return data.labels[data.labels.length - 1];
			return "";
		};
	}

	// tooltip disabled by default, each chart should implement its own with alignedTooltip
	options.plugins.tooltip = {
		enabled: false,
	};
};

const applyOptionsByChartId = (id, options, data) => {
	const defs = {
		cost_per_proof_chart: () => costPerProofCustomOptions(options, data),
	};

	return defs[id] ? defs[id]() : {};
};

export default {
	mounted() {
		this.initChart();
		window.addEventListener("theme-changed", this.reinitChart.bind(this));
	},

	updated() {
		this.reinitChart();
	},

	destroyed() {
		if (this.chart) {
			this.chart.destroy();
		}

		window.removeEventListener(
			"theme-changed",
			this.reinitChart.bind(this)
		);
	},

	initChart() {
		const ctx = this.el;
		const type = this.el.dataset.chartType;
		const data = JSON.parse(this.el.dataset.chartData);
		const options = JSON.parse(this.el.dataset.chartOptions);
		const chartId = this.el.id;

		applyCommonChartOptions(options, data);
		applyOptionsByChartId(chartId, options, data);

		this.chart = new Chart(ctx, {
			type,
			data,
			options,
		});
	},

	reinitChart() {
		if (this.chart) {
			this.chart.destroy();
		}
		this.initChart();
	},
};
