import { KecpMessageType } from './enums';

type RTCIceCredentialType = 'password';

interface RTCIceServer {
    credential?: string;
    credentialType?: RTCIceCredentialType;
    urls: string | string[];
    username?: string;
}

type Room = {
    roomID: string
    errorText: string
}

type CreateRoomResponse = {
    room_id: string
}

type ErrResponse = {
    status: string
    code: string
    error: string
}

type KecpConnectionOption = {
    roomID: string
    name: string
    iceServers?: RTCIceServer[]
}

type KecpInternalConnectionOption = {
    websocketURL: string
    roomID: string
    name: string
    clientKey: string
    iceServers: RTCIceServer[]
}

type KecpMessage = {
    type: KecpMessageType
    name?: string
    target?: string
    payload: any
}

export type {
  RTCIceServer,
  Room,
  CreateRoomResponse,
  ErrResponse,
  KecpConnectionOption,
  KecpInternalConnectionOption,
  KecpMessage,
};
