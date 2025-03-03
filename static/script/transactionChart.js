(function transactionChart() {
  const bugReportData = JSON.parse(
    document.getElementById("transactionData").textContent,
  );
  const ctx = document.getElementById("transactionChart");

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

  const data = bugReportData.daily.map((day) => day.count);

  new Chart(ctx, {
    type: "line",
    data: {
      labels: labels,
      datasets: [
        {
          data: data,
          tension: 0.3,
          fill: true,
          backgroundColor: "oklch(0.527 0.154 150.069/0.4)",
          borderColor: "oklch(0.527 0.154 150.069)",
        },
      ],
    },
    options: {
      scales: {
        y: {
          beginAtZero: true,
        },
      },
      plugins: {
        tooltip: {
          intersect: false,
        },
        legend: {
          display: false,
          labels: {
            boxWidth: 0,
          },
        },
      },
    },
    plugins: [],
  });
})();
