import adapter from 'webrtc-adapter';

class Browser {
  private supportsUnifiedPlan: boolean | undefined = adapter.browserDetails.supportsUnifiedPlan;

  get supportsWebrtc(): boolean {
    return !!this.supportsUnifiedPlan;
  }
}

const browser = new Browser();
export default browser;
