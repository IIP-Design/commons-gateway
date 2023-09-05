import { parseISO } from "date-fns";

/**
 * Decode JWT authentication token
 *
 * @param token string
 * @returns Object
 */
export const decodeJwt = (token:string) => {
  const base64Url = token?.split('.')[1];
  const base64 = base64Url?.replace(/-/g, '+')?.replace(/_/g, '/');
  const jsonPayload = decodeURIComponent(window.atob(base64).split('').map((c) => {
    const payload = c.charCodeAt(0).toString(16).slice(-2);
    return `%${payload}`;
  }).join(''));
  return JSON.parse(jsonPayload);
};

export const tokenExpiration = (token: string) => {
  const { exp } = decodeJwt( token );
  const dt = parseISO( exp );
  return dt.valueOf() / 1000;
}