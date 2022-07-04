import type KecpConnection from './connection';
import { KecpEventType, KecpMessageType } from './enums';
import { RTCIceServer } from './types';

export default class Peer {
  protected peerConnection: RTCPeerConnection;

  protected kecpConnection: KecpConnection;

  protected target: string;

  constructor(
    kecpConnection: KecpConnection,
    iceServers: RTCIceServer[],
    target: string,
  ) {
    this.kecpConnection = kecpConnection;
    this.peerConnection = new RTCPeerConnection({ iceServers });
    this.target = target;
    this.peerConnection.onicecandidate = (event) => this.handleICECandidateEvent(event);
    this.peerConnection.onconnectionstatechange = () => this.handleICEConnectionStateChangeEvent();
    this.peerConnection.onsignalingstatechange = () => this.handleSignalingStateChangeEvent();
    this.peerConnection.onnegotiationneeded = () => this.handleNegotiationNeededEvent();
    this.kecpConnection.on(
      KecpEventType.NewIceCandidate,
      (event) => this.handleNewICECandidateMsg(event),
    );
    this.kecpConnection.on(
      KecpEventType.VideoAnswer,
      (event) => this.handleVideoAnswerMsg(event),
    );
  }

  private async handleNewICECandidateMsg(event: CustomEvent) {
    const candidate = new RTCIceCandidate(event.detail.payload);
    try {
      await this.peerConnection.addIceCandidate(candidate);
    } catch (err) {
      this.peerConnectionClose();
    }
  }

  private async handleVideoAnswerMsg(event: CustomEvent) {
    const desc = new RTCSessionDescription(event.detail.payload);
    await this.peerConnection.setRemoteDescription(desc);
  }

  private handleICECandidateEvent(event: RTCPeerConnectionIceEvent) {
    if (event.candidate) {
      this.send(JSON.stringify({
        type: KecpMessageType.NewIceCandidate,
        name: this.kecpConnection.getName(),
        target: this.target,
        payload: event.candidate,
      }));
    }
  }

  private handleICEConnectionStateChangeEvent() {
    switch (this.peerConnection.iceConnectionState) {
      case 'closed':
      case 'failed':
      case 'disconnected':
        this.peerConnectionClose();
        break;
      default:
    }
  }

  private handleSignalingStateChangeEvent() {
    switch (this.peerConnection.signalingState) {
      case 'closed':
        this.peerConnectionClose();
        break;
      default:
    }
  }

  private async handleNegotiationNeededEvent() {
    const offer = await this.peerConnection.createOffer();

    if (this.peerConnection.signalingState !== 'stable') {
      return;
    }

    await this.peerConnection.setLocalDescription(offer);

    this.send(JSON.stringify({
      name: this.kecpConnection.getName(),
      target: this.target,
      type: 'video-offer',
      payload: this.peerConnection.localDescription,
    }));
  }

  protected send(msg: string) {
    this.kecpConnection.send(msg);
  }

  protected async peerConnectionClose() {
    this.peerConnection.ontrack = null;
    this.peerConnection.onicecandidate = null;
    this.peerConnection.oniceconnectionstatechange = null;
    this.peerConnection.onsignalingstatechange = null;
    this.peerConnection.onicegatheringstatechange = null;
    this.peerConnection.onnegotiationneeded = null;
    this.peerConnection.getTransceivers().forEach((track) => {
      track.stop();
    });
    this.peerConnection.close();
  }

  // eslint-disable-next-line no-unused-vars
  setHandleTrackEvent(callback: (event: RTCTrackEvent) => void) {
    this.peerConnection.ontrack = (event) => callback(event);
  }

  getTargetName(): string {
    return this.target;
  }
}
