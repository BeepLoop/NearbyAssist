const image = document.getElementById("image-viewer");
const controls = document.getElementById("image-controls").children;

for (const control of controls) {
    control.addEventListener("click", async () => {
        const btn = control.attributes["data-button"].value;
        const imageUrl = control.attributes["data-image"].value;

        image.setAttribute("src", imageUrl);
        image.setAttribute("alt", btn);
    })
}
