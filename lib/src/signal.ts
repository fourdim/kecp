import axios, { AxiosError } from 'axios';
import KecpRoom from './room';
import { genCryptoKey } from './helper';
import type {
  CreateRoomResponse, ErrResponse, KecpRoomOption,
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

  createRoom(): Promise<string> {
    return axios.post(this.roomsEndPoint, {
      client_key: this.clientKey,
    }, { timeout: 8000 }).then((res) => {
      const createRoomResponse = res.data as CreateRoomResponse;
      return createRoomResponse.room_id;
    }).catch((err: AxiosError) => {
      if (err.response) {
        const errResponse = err.response.data as ErrResponse;
        if (errResponse) {
          throw errResponse.error;
        }
        throw err.response.statusText;
      }
      throw err.message;
    });
  }

  getRoom(option: KecpRoomOption): KecpRoom {
    return new KecpRoom({
      websocketURL: this.clientsEndPoint,
      roomID: option.roomID,
      clientKey: this.clientKey,
      iceServers: option.iceServers ? option.iceServers : [{
        urls: 'stun:stun.stunprotocol.org',
      }],
    });
  }
}
