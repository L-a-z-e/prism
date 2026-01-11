import { Client } from '@stomp/stompjs';
import { onUnmounted } from 'vue';

export const useWebSocket = (onConnect?: () => void) => {
  const client = new Client({
    brokerURL: 'ws://localhost:8085/ws',
    onConnect: () => {
      console.log('Connected to WebSocket');
      if (onConnect) onConnect();
    },
    onStompError: (frame) => {
      console.error('Broker reported error: ' + frame.headers['message']);
      console.error('Additional details: ' + frame.body);
    },
  });

  const connect = () => {
    client.activate();
  };

  const disconnect = () => {
    client.deactivate();
  };

  const subscribe = (destination: string, callback: (message: any) => void) => {
    return client.subscribe(destination, (message) => {
      callback(JSON.parse(message.body));
    });
  };

  onUnmounted(() => {
    disconnect();
  });

  return { client, connect, disconnect, subscribe };
};
