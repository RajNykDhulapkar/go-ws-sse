import { WebSocket } from "websockets";

// Connect to the WebSocket server on the specified route
const ws = new WebSocket("ws://localhost:8080/api/ws");

ws.on("open", function open() {
  console.log("Connected to the server.");
  // Send a message to the server once the connection is opened.
  ws.send("Hello, server!");
});

ws.on("message", function incoming(data) {
  // Log messages received from the server
  console.log("Message from server:", data);
});

ws.on("error", function handleError(error) {
  // Handle any errors that occur
  console.error("WebSocket error:", error);
});

ws.on("close", function handleClose(code, reason) {
  // Handle the connection closing
  console.log(`WebSocket closed with code: ${code} and reason: ${reason}`);
});
