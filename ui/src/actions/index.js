import { API_BASE_URL } from '../config';
export default function apiFetch(endpoint, options = {}) {
  options.headers = {
    // 'X-JWT': cookie.load('token')
  };
  return fetch(`${API_BASE_URL}/${endpoint}`, options);
}
