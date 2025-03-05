(function userDataChart() {
  const userData = JSON.parse(document.getElementById("userData").textContent);

  const { total, verified, expert, restricted } = userData;
  const unverified = total - verified;

  const data = [
    {
      type: "pie",
      hole: 0.6,
      labels: ["unverified", "restricted", "total"],
      values: [unverified, restricted, total - restricted - unverified],
      name: "verified vs unverified",
      marker: {
        colors: [
          "rgba(234, 168, 106, 1)",
          "rgba(255, 60, 40, 1)",
          "rgba(96, 120, 214, 0.75)",
        ],
      },
      automargin: true,
    },
  ];

  const layout = {
    showlegend: false,
    margin: { t: 0, b: 0, l: 40, r: 40 },
  };

  const options = {
    displayModeBar: false,
    responsive: true,
  };

  Plotly.newPlot("userDataChart", data, layout, options);
})();
