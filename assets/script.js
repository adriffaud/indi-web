window.addEventListener("load", () => {
  sse_client();
});

function sse_client() {
  let sse = new EventSource("http://localhost:8080/sse");

  sse.onopen = () => console.log("SSE open");

  sse.onmessage = (evt) => {
    console.log(evt.data);
  };

  sse.onclose = () => console.log("SSE closed");
}

function ws_client() {
  let ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = () => console.log("WS open");

  ws.onmessage = (evt) => {
    console.log(`${new Date()}: ${evt.data}`);
  };

  ws.onclose = () => console.log("WS closed");
}
