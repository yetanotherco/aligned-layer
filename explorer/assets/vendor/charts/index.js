import Chart from "chart.js/auto";
import { costPerProofCustomOptions } from "./cost_per_proof";

const applyCommonChartOptions = (options) => {
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

		applyCommonChartOptions(options);
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
