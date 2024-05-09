const grpc = require('grpc-web');
const { Empty, Ping } = require('./pingpong_pb.js');
const { PingPongServiceClient } = require('./pingpong_grpc_web_pb.js');

const serverUrl = process.env.GRPC_SERVER_URL || 'http://localhost:8080';
const client = new PingPongServiceClient(serverUrl, null, null);

function listenForPongs() {
  const emptyRequest = new Empty();
  const stream = client.receivePongStream(emptyRequest);

  stream.on('data', function(response) {
    const message = document.createElement('li');
    message.textContent = `Received Pong: ${response.getMessage()}`;
    document.getElementById('messages').appendChild(message);
  });

  stream.on('status', function(status) {
    if (status.code !== grpc.StatusCode.OK) {
      console.error('Stream status error:', status.details);
    }
  });

  stream.on('end', function() {
    console.log('Stream ended');
  });
}

function sendPing() {
  const pingRequest = new Ping();
  pingRequest.setMessage("Ping");
  const stream = client.pingAndStreamPong(pingRequest);

  stream.on('data', function(response) {
    const message = document.createElement('li');
    message.textContent = `Received Pong for Ping: ${response.getMessage()}`;
    document.getElementById('messages').appendChild(message);
  });

  stream.on('status', function(status) {
    if (status.code !== grpc.StatusCode.OK) {
      console.error('Error status:', status.details);
    }
  });

  stream.on('end', function() {
    console.log('Stream from Ping ended');
  });
}

// ウェブページがロードされたときにPongを受け取るように設定
window.addEventListener('load', listenForPongs);
window.sendPing = sendPing;
