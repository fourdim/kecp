import adapter from 'webrtc-adapter';

class Browser {
  private supportsUnifiedPlan: boolean | undefined = adapter.browserDetails.supportsUnifiedPlan;

  private browserBrand: string = adapter.browserDetails.browser;

  get supportsWebrtc(): boolean {
    return !!this.supportsUnifiedPlan;
  }

  get browser(): string {
    return this.browserBrand;
  }
}

const browser = new Browser();
export default browser;
