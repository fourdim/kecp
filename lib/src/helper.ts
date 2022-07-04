export function combineURLs(baseURL: string, relativeURL: string) {
  return relativeURL
    ? `${baseURL.replace(/\/+$/, '')}/${relativeURL.replace(/^\/+/, '')}`
    : baseURL;
}

export function genCryptoKey(): string {
  const b = new Uint8Array(48);
  window.crypto.getRandomValues(b);
  let binary = '';
  const len = b.byteLength;
  for (let i = 0; i < len; i += 1) {
    binary += String.fromCharCode(b[i]);
  }
  return window.btoa(binary).replace('+', '-').replace('/', '_');
}
