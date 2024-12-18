import Chart from "chart.js/auto";
import { costPerProofCustomOptions } from "./cost_per_proof";

const cssvar = (name) =>
	getComputedStyle(document.documentElement)
		.getPropertyValue(`--${name}`)
		.trim();
const cssColor = (name, opacity = 1) => `hsl(${cssvar(name)} / ${opacity})`;

const applyCommonChartOptions = (options, tooltip) => {
	options.plugins.tooltip = {
		displayColors: false,
		backgroundColor: cssColor("card"),
		borderColor: cssColor("foreground", 0.2),
		borderWidth: 1,
		cornerRadius: 8,
		padding: {
			x: 20,
			y: 15,
		},
		titleColor: cssColor("foreground"),
		titleAlign: "center",
		titleMarginBottom: 10,
		titleFont: {
			size: 12,
		},
		bodyColor: cssColor("foreground"),
		bodyAlign: "left",
		bodySpacing: 5,
		footerColor: cssColor("muted-foreground"),
		footerAlign: "center",
		callbacks: {
			title: () => tooltip?.title,
			label: () => "",
			afterBody: (item) => {
				const value = item[0].formattedValue;
				const label = item[0].label;
				return tooltip?.body
					.replace("{{value}}", value)
					.replace("{{label}}", label);
			},
		},
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
		const tooltip = JSON.parse(this.el.dataset.chartTooltip);
		const chartId = this.el.id;

		applyCommonChartOptions(options, tooltip);
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
