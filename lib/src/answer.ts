import type KecpRoom from './room';
import { KecpMessageType } from './enums';
import Peer from './peer';
import type { KecpMessage, RTCIceServer } from './types';

export default class AnswerPeer extends Peer {
  private sdp: any;

  constructor(offer: KecpMessage, room: KecpRoom, iceServers: RTCIceServer[]) {
    super(room, iceServers, offer.name!);
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
      name: this.kecpRoom.getSelfName(),
      target: this.target,
      payload: this.peerConnection.localDescription,
    }));
  }

  async preAnswer() {
    const desc = new RTCSessionDescription(this.sdp);
    await this.peerConnection.setRemoteDescription(desc);
    const preAnswerDesc = await this.peerConnection.createAnswer();
    preAnswerDesc.sdp = preAnswerDesc.sdp!.replaceAll(/a=recvonly/g, 'a=inactive');
    preAnswerDesc.type = 'pranswer';
    await this.peerConnection.setLocalDescription(preAnswerDesc);
    this.send(JSON.stringify({
      type: KecpMessageType.VideoAnswer,
      name: this.kecpRoom.getSelfName(),
      target: this.target,
      payload: preAnswerDesc,
    }));
  }

  async preAnswerConfirm() {
    const answerDesc = await this.peerConnection.createAnswer();
    answerDesc.sdp = answerDesc.sdp!.replaceAll(/a=inactive/g, 'a=recvonly');
    answerDesc.type = 'answer';
    await this.peerConnection.setLocalDescription(answerDesc);
    this.send(JSON.stringify({
      type: KecpMessageType.VideoAnswer,
      name: this.kecpRoom.getSelfName(),
      target: this.target,
      payload: answerDesc,
    }));
  }
}
