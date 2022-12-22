import Peer from './video-peer';

export default class VideoOfferPeer extends Peer {
  addStream(stream: MediaStream) {
    stream.getTracks().forEach(
      (track) => this.peerConnection.addTrack(track, stream),
    );
  }

  setBandwidth(bandWidth: number) {
    this.bandWidth = bandWidth;
  }
}
