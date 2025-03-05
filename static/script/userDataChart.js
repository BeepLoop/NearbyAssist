(function userDataChart() {
  const userData = JSON.parse(document.getElementById("userData").textContent);

  const labels = ["unverified", "experts"];
  const values = [userData.total - userData.verified, userData.expert];

  const data = [
    {
      type: "pie",
      hole: 0.5,
      labels: labels,
      values: values,
      automargin: true,
    },
  ];

  const layout = {
    showlegend: false,
    margin: { t: 0, b: 0, l: 0, r: 0 },
  };

  const options = {
    displayModeBar: false,
    responsive: true,
  };

  Plotly.newPlot("userDataChart", data, layout, options);
})();
