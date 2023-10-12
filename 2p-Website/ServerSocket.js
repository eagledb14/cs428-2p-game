import { WebSocketServer } from 'ws';

const wss = new WebSocketServer({ port: 3001 });

wss.on('connection', (ws) => {
  ws.on('error', console.error);

  ws.on('message', (data) => {
    console.log('received: %s', data);
  });

  ws.on('open', (event) => {
    console.log('connected', event)
  })

  ws.send(JSON.stringify({message: 'welcome to the connection'}));
});