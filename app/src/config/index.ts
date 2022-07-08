export const kecpEndpoint = import.meta.env.PROD
  ? `${window.location.protocol}//${window.location.host}/api/kecp`
  : "http://127.0.0.1:8090/api/kecp/";
