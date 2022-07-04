import type KecpConnection from './connection';
import { KecpMessageType } from './enums';
import Peer from './peer';
import { KecpMessage, RTCIceServer } from './types';

export default class AnswerPeer extends Peer {
  private sdp: any;

  constructor(offer: KecpMessage, kecpConnection: KecpConnection, iceServers: RTCIceServer[]) {
    super(kecpConnection, iceServers, offer.name!);
    this.sdp = offer.payload;
  }

  async answer() {
    const desc = new RTCSessionDescription(this.sdp);
    if (this.peerConnection.signalingState !== 'stable') {
      // Set the local and remove descriptions for rollback; don't proceed
      // until both return.
      await Promise.all([
        this.peerConnection.setLocalDescription({ type: 'rollback' }),
        this.peerConnection.setRemoteDescription(desc),
      ]);
    } else {
      await this.peerConnection.setRemoteDescription(desc);
    }
    await this.peerConnection.setLocalDescription(await this.peerConnection.createAnswer());
    this.send(JSON.stringify({
      type: KecpMessageType.VideoAnswer,
      name: this.kecpConnection.getName(),
      target: this.target,
      payload: this.peerConnection.localDescription,
    }));
  }
}
