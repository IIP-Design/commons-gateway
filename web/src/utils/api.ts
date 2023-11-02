// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { accessToken } from '../stores/current-user';
import { logout } from './login';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
export type TMethods = 'DELETE' | 'GET' | 'POST' | 'PUT';

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
const API_ENDPOINT = import.meta.env.PUBLIC_SERVERLESS_URL; // The URL for the API Gateway.

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
export const constructUrl = ( endpoint: string ) => `${API_ENDPOINT}/${endpoint}`;

// ////////////////////////////////////////////////////////////////////////////
// API Functions
// ////////////////////////////////////////////////////////////////////////////

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const buildHeaders = ( token: string, body: Nullable<Record<string, any>> ): HeadersInit => {
  const headers: HeadersInit = {};

  if ( body ) {
    headers['Content-Type'] = 'application/json';
  }

  if ( token ) {
    headers.authorization = `Bearer ${token}`;
  }

  return headers;
};

/**
 * Helper function to consistently construct the API requests.
 * @param endpoint The API endpoint for the function in question (without a leading slash)
 * @param body The data to be sent to the API.
 * @param method The HTTP request method (if not provided defaults to POST).
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const buildQuery = async ( endpoint: string, body: Nullable<Record<string, any>>, method?: TMethods ) => {
  let opts = {
    headers: buildHeaders( accessToken.get(), body ),
    method: method || 'POST',
  } as RequestInit;

  if ( body !== null ) {
    opts = {
      ...opts,
      body: JSON.stringify( body as BodyInit ),
    };
  }

  const response = await fetch( constructUrl( endpoint ), opts );

  // 401 means the token is expired, so log out user for UX reasons
  if ( response.status === 401 ) {
    logout();
  }

  return response;
};
