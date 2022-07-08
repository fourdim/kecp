import AnswerPeer from './answer';
import OfferPeer from './offer';
import { KecpEventType, KecpMessageType } from './enums';
import type { KecpInternalConnectionOption, KecpMessage, RTCIceServer } from './types';

interface KecpConnectionCustomEventCallback {
  // eslint-disable-next-line no-unused-vars
  (evt: CustomEvent): void
}

interface KecpConnectionEventCallback {
  // eslint-disable-next-line no-unused-vars
  (evt: Event): void
}

export default class KecpConnection extends EventTarget {
  private ws: WebSocket;

  private roomID: string;

  private name: string;

  private clientKey: string;

  private iceServers: RTCIceServer[];

  private userList: string[];

  constructor(option: KecpInternalConnectionOption) {
    super();
    this.ws = new WebSocket(option.websocketURL);
    this.roomID = option.roomID;
    this.name = option.name;
    this.clientKey = option.clientKey;
    this.iceServers = option.iceServers;
    this.userList = [];
    // Avoid this being changed
    this.ws.onopen = () => this.onOpenHandler();
    this.ws.onmessage = (event) => this.onMessageHandler(event);
  }

  get isOpen(): boolean {
    return !!this.ws && this.ws.readyState === 1;
  }

  private onOpenHandler() {
    this.ws.send(JSON.stringify({
      room_id: this.roomID,
      name: this.name,
      client_key: this.clientKey,
    }));
    this.dispatchEvent(new CustomEvent(KecpEventType.Open));
  }

  private onMessageHandler(event: MessageEvent) {
    const message = JSON.parse(event.data) as KecpMessage;
    let messageEvent;
    switch (message.type) {
      case KecpMessageType.Error:
        messageEvent = new CustomEvent(KecpEventType.Error, { detail: message });
        this.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.VideoOffer:
        messageEvent = new CustomEvent(KecpEventType.VideoOffer, {
          detail: new AnswerPeer(message, this, this.iceServers),
        });
        this.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.VideoAnswer:
        messageEvent = new CustomEvent(KecpEventType.VideoAnswer, { detail: message });
        this.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.NewIceCandidate:
        messageEvent = new CustomEvent(KecpEventType.NewIceCandidate, { detail: message });
        this.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.Chat:
        messageEvent = new CustomEvent(KecpEventType.Chat, { detail: message });
        this.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.List:
        this.userList = message.payload;
        messageEvent = new CustomEvent(KecpEventType.UserListInit, { detail: message.payload });
        this.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.Join:
        this.userList.push(message.payload);
        messageEvent = new CustomEvent(KecpEventType.UserJoin, { detail: message.payload });
        this.dispatchEvent(messageEvent);
        break;
      case KecpMessageType.Leave:
        if (this.userList.indexOf(message.payload) !== -1) {
          this.userList.splice(this.userList.indexOf(message.payload));
        }
        messageEvent = new CustomEvent(KecpEventType.UserLeave, { detail: message.payload });
        this.dispatchEvent(messageEvent);
        break;
      default:
    }
  }

  newOffer(target: string): OfferPeer | undefined {
    if (this.userList.includes(target)) {
      return new OfferPeer(this, this.iceServers, target);
    }
    return undefined;
  }

  getName(): string {
    return this.name;
  }

  on(event: KecpEventType, callback: KecpConnectionCustomEventCallback) {
    this.addEventListener(event, ((e: CustomEvent) => callback(e)) as KecpConnectionEventCallback);
  }

  send(msg: string) {
    if (!this.isOpen) {
      return;
    }
    this.ws.send(msg);
  }

  close() {
    if (!this.isOpen) {
      return;
    }
    this.ws.onopen = null;
    this.ws.onmessage = null;
    this.ws.close();
  }
}
