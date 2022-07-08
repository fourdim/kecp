import axios, { AxiosError } from 'axios';
import KecpConnection from './connection';
import { genCryptoKey } from './helper';
import type {
  CreateRoomResponse, ErrResponse, KecpConnectionOption, Room,
} from './types';

export default class KecpSignal {
  private tls: boolean = window.location.protocol === 'https:';

  private roomsEndPoint: string;

  private clientsEndPoint: string;

  private clientKey: string;

  constructor(endPoint: string) {
    this.clientKey = genCryptoKey();
    const httpEndPoint = new URL(endPoint);
    const wsEndPoint = new URL(endPoint);
    httpEndPoint.protocol = this.tls ? 'https:' : 'http:';
    wsEndPoint.protocol = this.tls ? 'wss:' : 'ws:';
    this.roomsEndPoint = httpEndPoint.toString();
    this.clientsEndPoint = wsEndPoint.toString();
  }

  newRoom(): Promise<Room> {
    return axios.post(this.roomsEndPoint, {
      client_key: this.clientKey,
    }, { timeout: 8000 }).then((res) => {
      const createRoomResponse = res.data as CreateRoomResponse;
      return {
        roomID: createRoomResponse.room_id,
        errorText: '',
      };
    }).catch((err: AxiosError) => {
      if (err.response) {
        const errResponse = err.response.data as ErrResponse;
        if (errResponse) {
          return {
            roomID: '',
            errorText: errResponse.error,
          };
        }
        return {
          roomID: '',
          errorText: err.response.statusText,
        };
      }
      return {
        roomID: '',
        errorText: err.message,
      };
    });
  }

  newConnection(option: KecpConnectionOption): KecpConnection {
    return new KecpConnection({
      websocketURL: this.clientsEndPoint,
      roomID: option.roomID,
      name: option.name,
      clientKey: this.clientKey,
      iceServers: option.iceServers ? option.iceServers : [{
        urls: 'stun:stun.stunprotocol.org',
      }],
    });
  }
}
