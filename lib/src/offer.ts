import Peer from './peer';

export default class OfferPeer extends Peer {
  addStream(stream: MediaStream) {
    stream.getTracks().forEach(
      (track) => this.peerConnection.addTrack(track, stream),
    );
  }
}
