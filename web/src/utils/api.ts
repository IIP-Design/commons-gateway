// The URL for the API Gateway.
const API_ENDPOINT = import.meta.env.PUBLIC_SERVERLESS_URL;

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
export type TActions = 'create' | 'confirm';
export type TMethods = 'GET' | 'POST';

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
export const constructUrl = ( endpoint: string ) => `${API_ENDPOINT}/${endpoint}`;

// ////////////////////////////////////////////////////////////////////////////
// API Functions
// ////////////////////////////////////////////////////////////////////////////

/**
 * Helper function to consistently construct the API requests.
 * @param endpoint The API endpoint for the function in question (without a leading slash)
 * @param body The data to be sent to the API.
 * @param method The HTTP request method (if not provided defaults to POST).
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const buildQuery = async ( endpoint: string, body: Nullable<Record<string, any>>, method?: TMethods ) => {
  let opts = {
    headers: {
      'Content-Type': 'application/json',
    },
    method: method || 'POST',
  } as RequestInit;

  if ( body !== null ) {
    opts = {
      ...opts,
      body: JSON.stringify( body as BodyInit ),
    };
  }

  return fetch( constructUrl( endpoint ), opts );
};
