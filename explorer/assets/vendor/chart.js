import Chart from "chart.js/auto";

export default {
	mounted() {
		const ctx = this.el;
		const type = this.el.dataset.chartType;
		const data = JSON.parse(this.el.dataset.chartData);
		const options = JSON.parse(this.el.dataset.chartOptions);
		console.log(data, options);

		new Chart(ctx, {
			type,
			data,
			options,
		});
	},
	updated() {},
};
