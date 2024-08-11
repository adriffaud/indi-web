const buttons = document.querySelectorAll("button");

for (let i = 0; i < buttons.length; i++) {
  buttons[i].addEventListener("click", async (evt) => {
    const response = await fetch("/indi/action", { method: "POST", body: JSON.stringify(evt.target.dataset) });
    console.log(response);
  });
}
