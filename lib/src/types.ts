import type { KecpMessageType } from './enums';

type RTCIceCredentialType = 'password';

interface RTCIceServer {
    credential?: string;
    credentialType?: RTCIceCredentialType;
    urls: string | string[];
    username?: string;
}

type CreateRoomResponse = {
    room_id: string
}

type ErrResponse = {
    status: string
    code: string
    error: string
}

type KecpRoomOption = {
    roomID: string
    iceServers?: RTCIceServer[]
}

type KecpRoomInternalOption = {
    websocketURL: string
    roomID: string
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
  CreateRoomResponse,
  ErrResponse,
  KecpRoomOption,
  KecpRoomInternalOption,
  KecpMessage,
};
