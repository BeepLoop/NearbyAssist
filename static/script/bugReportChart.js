const data = JSON.parse(document.getElementById("bugReportData").textContent);
const ctx = document.getElementById("bugReportChart");

const labels = [1, 2, 3, 4, 5, 6, 7];

new Chart(ctx, {
  type: "line",
  data: {
    labels: labels,
    datasets: [
      {
        label: "# of reports",
        data: [],
        borderWidth: 1,
        fill: true,
        tension: 0.3,
      },
    ],
  },
  options: {
    scales: {
      y: {
        suggestedMin: 0,
        suggestedMax: 10,
      },
    },
  },
  plugins: [],
});
