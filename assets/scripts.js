const buttons = document.querySelectorAll("button");

for (let i = 0; i < buttons.length; i++) {
  buttons[i].addEventListener("click", async (evt) => await fetch("/indi/action", { method: "POST", body: JSON.stringify(evt.target.dataset) }));
}

window.addEventListener("load", (evt) => {
  const output = document.getElementById("output");

  let ws;

  const print = function(message) {
    const d = document.createElement("div");
    d.textContent = message;
    output.appendChild(d);
    output.scroll(0, output.scrollHeight);
  };

  ws = new WebSocket("ws://localhost:8080/ws");
  ws.onopen = () => {
    print("OPEN");
  };
  ws.onclose = () => {
    print("CLOSE");
    ws = null;
  }
  ws.onmessage = (evt) => {
    print("RESPONSE: " + evt.data);
  }
  ws.onerror = (evt) => {
    print("ERROR: " + evt.data);
  }
});
