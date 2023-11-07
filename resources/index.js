import { WebSocket } from "websockets";
import EventSource from "eventsource";

const NUM_SENDER_CLIENTS = 5; // Define the number of sender clients
const NUM_RECEIVER_PER_SENDER = 3; // Define the number of receivers per sender

// Function to generate a random room ID
function generateRandomRoomId() {
  return Math.random().toString(36).substring(2, 10);
}

// Function to initialize a WebSocket connection
function initializeWebSocket(roomId) {
  const ws = new WebSocket(`ws://localhost:8080/api/ws?roomId=${roomId}`);

  ws.on("open", () => {
    console.log(`WebSocket client connected to roomId: ${roomId}`);
    setInterval(() => {
      const message =
        new Date().toLocaleTimeString() + " " + roomId + " " + Math.random();
      ws.send(message);
    }, 1000);
  });

  ws.on("message", (data) => {
    console.log(`Message from server: ${data}`);
  });

  ws.on("error", (error) => {
    console.error(`WebSocket error: ${error.message}`);
  });
}

// Function to initialize an SSE connection
function initializeSSE(roomId) {
  const eventSource = new EventSource(
    `http://localhost:8080/api/sse?roomId=${roomId}`
  );

  eventSource.onmessage = (event) => {
    console.log(`SSE Message for roomId ${roomId}:`, event.data);
  };

  eventSource.onerror = (error) => {
    console.error(`SSE error for roomId ${roomId}:`, error);
  };

  eventSource.onopen = () => {
    console.log(`SSE opened for roomId ${roomId}`);
  };
}

// Generate room IDs and initialize clients
const roomIds = Array.from(
  { length: NUM_SENDER_CLIENTS },
  generateRandomRoomId
);

// Initialize WebSocket clients
roomIds.forEach(initializeWebSocket);

// Initialize SSE clients for each room
roomIds.forEach((roomId) => {
  Array.from({ length: NUM_RECEIVER_PER_SENDER }).forEach(() => {
    initializeSSE(roomId);
  });
});
