const {Ping, Pong} = require('./pingpong_pb.js');
const {PingPongServiceClient} = require('./pingpong_grpc_web_pb.js');

const client = new PingPongServiceClient('http://localhost:8080', null, null);

function startPingPong() {
  const stream = client.streamPingPong();
  stream.on('data', function (response) {
    const message = document.createElement('li');
    message.textContent = `Received: ${response.getMessage()}`;
    document.getElementById('messages').appendChild(message);
  });

  stream.on('status', function (status) {
    if (status.code !== grpc.Code.OK) {
      console.log('Error status: ' + status.details);
    }
  });

  // Send pings continuously
  setInterval(() => {
    const request = new Ping();
    request.setMessage("Ping");
    stream.write(request);
  }, 1000);

  stream.on('end', function (end) {
    // stream end signal
  });
}
