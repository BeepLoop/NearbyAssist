(function transactionChart() {
  const bugReportData = JSON.parse(
    document.getElementById("transactionData").textContent,
  );

  const months = [
    "January",
    "February",
    "March",
    "April",
    "May",
    "June",
    "July",
    "August",
    "September",
    "October",
    "November",
    "December",
  ];

  const labels = bugReportData.daily.map((day) => {
    const date = new Date(day.date);
    return `${months[date.getMonth()]} ${date.getDate()}`;
  });

  const values = bugReportData.daily.map((day) => day.count);

  const DEFAULT_MAX_RANGE = 10;
  const max = Math.max(...values);
  const range = [0, max > DEFAULT_MAX_RANGE ? max : DEFAULT_MAX_RANGE];

  const data = [
    {
      type: "bar",
      x: labels,
      y: values,
      marker: {
        line: {
          color: "rgb(10, 10, 255)",
        },
        opacity: 0.5,
      },
      width: 0.5,
      automargin: true,
    },
  ];

  const layout = {
    showlegend: false,
    xaxis: {
      tickangle: -45,
      tickwidth: 4,
      zeroline: true,
      showgrid: true,
      showline: true,
    },
    yaxis: {
      gridwidth: 2,
      tickwidth: 4,
      range: range,
    },
    bargap: 0.05,
  };

  const options = {
    displayModeBar: false,
    responsive: true,
  };

  Plotly.newPlot("transactionChart", data, layout, options);
})();
