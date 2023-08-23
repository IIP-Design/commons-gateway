// The URL for the API Gateway.
const API_ENDPOINT = import.meta.env.PUBLIC_SERVERLESS_URL;

type TActions = 'create' | 'confirm';
type TMethods = 'GET' | 'POST';

/** The values that can be sent to the server. */
interface IFetchBody {
  action?: TActions
  hash?: string
  username?: string
}

/**
 * Helper function to consistently construct the API requests.
 * @param endpoint The API endpoint for the function in question (without a leading slash)
 * @param body The data to be sent to the API.
 * @param method The HTTP request method.
 */
const buildQuery = async ( endpoint: string, body: IFetchBody, method?: TMethods ) => (
  fetch( `${API_ENDPOINT}/${endpoint}`, {
    body: JSON.stringify( body ),
    headers: {
      'Content-Type': 'application/json',
    },
    method: method || 'POST',
  } )
);

/**
 * Retrieves the salt value used to hash the user's password.
 * @param username The name of the user to look up.
 * @returns The salt value (if the user exits).
 */
export const passTheSalt = async ( username: string ) => {
  const response = await buildQuery( 'creds/salt', { username } );

  const { salt } = await response.json();

  return salt;
};

/**
 * Send the locally generated password hash to the server to authenticate user and request access.
 * @param action Whether to initiate a authenticated session or confirm an existing session.
 * @param hash The locally generated password hash.
 * @param username The email of the user attempting to log in.
 */
export const submitHash = async ( action: TActions, hash: string, username: string ) => {
  const response = await buildQuery( 'creds', {
    action,
    hash,
    username,
  } );

  const res = await response.json();

  console.log( res );
};
