import type KecpRoom from './room';
import { KecpEventType, KecpMessageType } from './enums';
import type { RTCIceServer } from './types';

export default class Peer {
  protected peerConnection: RTCPeerConnection;

  protected kecpRoom: KecpRoom;

  protected target: string;

  protected bandWidth: number | undefined;

  constructor(
    room: KecpRoom,
    iceServers: RTCIceServer[],
    target: string,
  ) {
    this.kecpRoom = room;
    this.peerConnection = new RTCPeerConnection({ iceServers });
    this.target = target;
    this.peerConnection.onicecandidate = (event) => this.handleICECandidateEvent(event);
    this.peerConnection.onconnectionstatechange = () => this.handleICEConnectionStateChangeEvent();
    this.peerConnection.onsignalingstatechange = () => this.handleSignalingStateChangeEvent();
    this.peerConnection.onnegotiationneeded = () => this.handleNegotiationNeededEvent();
    this.kecpRoom.on(
      KecpEventType.NewIceCandidate,
      (event) => this.handleNewICECandidateMsg(event),
    );
    this.kecpRoom.on(
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
    if (desc.type === 'answer' && this.bandWidth !== undefined) {
      const sender = this.peerConnection.getSenders()[0];
      const parameters = sender.getParameters();
      if (!parameters.encodings) {
        parameters.encodings = [{}];
      }
      parameters.encodings[0].maxBitrate = this.bandWidth * 1000;
      sender.setParameters(parameters).catch();
    }
  }

  private handleICECandidateEvent(event: RTCPeerConnectionIceEvent) {
    if (event.candidate) {
      this.send(JSON.stringify({
        type: KecpMessageType.NewIceCandidate,
        name: this.kecpRoom.getSelfName(),
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
      name: this.kecpRoom.getSelfName(),
      target: this.target,
      type: KecpMessageType.VideoOffer,
      payload: this.peerConnection.localDescription,
    }));
  }

  protected send(msg: string) {
    this.kecpRoom.send(msg);
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

  close() {
    this.peerConnectionClose();
  }
}
