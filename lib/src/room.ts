import AnswerPeer from './video-answer';
import VideoOfferPeer from './video-offer';
import { KecpEventType, KecpMessageType } from './enums';
import type { KecpRoomInternalOption, KecpMessage, RTCIceServer } from './types';

interface KecpRoomCustomEventCallback {
  // eslint-disable-next-line no-unused-vars
  (evt: CustomEvent): void
}

interface KecpRoomEventCallback {
  // eslint-disable-next-line no-unused-vars
  (evt: Event): void
}

export default class KecpRoom {
  private et: EventTarget;

  private ws: WebSocket | undefined;

  private wsURL: string;

  private roomID: string;

  private name: string;

  private clientKey: string;

  private iceServers: RTCIceServer[];

  private userList: string[];

  constructor(option: KecpRoomInternalOption) {
    this.et = new EventTarget();
    this.wsURL = option.websocketURL;
    this.roomID = option.roomID;
    this.name = '';
    this.clientKey = option.clientKey;
    this.iceServers = option.iceServers;
    this.userList = [];
  }

  private onOpenHandler() {
    if (!this.isOpen()) {
      return;
    }
    this.ws!.send(JSON.stringify({
      room_id: this.roomID,
      name: this.name,
      client_key: this.clientKey,
    }));
    this.et.dispatchEvent(new CustomEvent(KecpEventType.Open));
  }

  private onMessageHandler(event: MessageEvent) {
    const message = JSON.parse(event.data) as KecpMessage;
    let messageEvent;
    switch (message.type) {
      case KecpMessageType.Error:
        messageEvent = new CustomEvent(KecpEventType.Error, { detail: message });
        this.et.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.VideoOffer:
        messageEvent = new CustomEvent(KecpEventType.VideoOffer, {
          detail: new AnswerPeer(message, this, this.iceServers),
        });
        this.et.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.VideoAnswer:
        messageEvent = new CustomEvent(KecpEventType.VideoAnswer, { detail: message });
        this.et.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.NewIceCandidate:
        messageEvent = new CustomEvent(KecpEventType.NewIceCandidate, { detail: message });
        this.et.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.Chat:
        messageEvent = new CustomEvent(KecpEventType.Chat, { detail: message });
        this.et.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.List:
        this.userList = message.payload;
        messageEvent = new CustomEvent(KecpEventType.UserListInit, { detail: message.payload });
        this.et.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.Join:
        this.userList.push(message.payload);
        messageEvent = new CustomEvent(KecpEventType.UserJoin, { detail: message.payload });
        this.et.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.Leave:
        if (this.userList.indexOf(message.payload) !== -1) {
          this.userList.splice(this.userList.indexOf(message.payload));
        }
        messageEvent = new CustomEvent(KecpEventType.UserLeave, { detail: message.payload });
        this.et.dispatchEvent(messageEvent);
        break;
      default:
    }
  }

  isOpen(): boolean {
    return !!this.ws && this.ws.readyState === 1;
  }

  connect(username: string) {
    this.disconnect();
    this.name = username;
    this.ws = new WebSocket(this.wsURL);
    // Avoid this being changed
    this.ws.onopen = () => this.onOpenHandler();
    this.ws.onmessage = (event) => this.onMessageHandler(event);
  }

  newVideoOffer(target: string): VideoOfferPeer | undefined {
    if (this.userList.includes(target)) {
      return new VideoOfferPeer(this, this.iceServers, target);
    }
    return undefined;
  }

  newDataOffer(target: string): VideoOfferPeer | undefined {
    if (this.userList.includes(target)) {
      return new VideoOfferPeer(this, this.iceServers, target);
    }
    return undefined;
  }

  getSelfName(): string {
    return this.name;
  }

  on(event: KecpEventType, callback: KecpRoomCustomEventCallback) {
    this.et.addEventListener(event, ((e: CustomEvent) => callback(e)) as KecpRoomEventCallback);
  }

  send(msg: string) {
    if (!this.isOpen()) {
      return;
    }
    this.ws!.send(msg);
  }

  disconnect() {
    if (!this.isOpen()) {
      return;
    }
    this.ws!.onopen = null;
    this.ws!.onmessage = null;
    this.ws!.close();
  }
}
