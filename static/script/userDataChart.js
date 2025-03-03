(function userDataChart() {
  const userData = JSON.parse(document.getElementById("userData").textContent);
  const ctx = document.getElementById("userDataChart");

  const labels = ["unverified", "experts"];
  const data = [userData.total - userData.verified, userData.expert];

  new Chart(ctx, {
    type: "doughnut",
    data: {
      labels: labels,
      datasets: [
        {
          data: data,
          backgroundColor: [
            "oklch(0.637 0.237 25.331/0.5)",
            "oklch(0.527 0.154 150.069/0.5)",
          ],
          hoverOffset: 4,
        },
      ],
    },
    options: {
      cutout: "75%",
      radius: "90%",
      plugins: {
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
