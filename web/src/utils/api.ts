// The URL for the API Gateway.
const API_ENDPOINT = import.meta.env.PUBLIC_SERVERLESS_URL;

type TActions = 'create' | 'confirm';
type TMethods = 'GET' | 'POST';

/** The values that can be sent to the server. */
interface IFetchBody {
  action?: TActions
  active?: boolean
  hash?: string
  team?: string
  teamName?: string
  username?: string
  inviter?: string,
  invitee?: {
    email?: string,
    givenName?: string,
    familyName?: string,
    team?: string,
  },
}

/**
 * Helper function to consistently construct the API requests.
 * @param endpoint The API endpoint for the function in question (without a leading slash)
 * @param body The data to be sent to the API.
 * @param method The HTTP request method.
 */
export const buildQuery = async ( endpoint: string, body: IFetchBody | null, method?: TMethods ) => {
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

  return (
    fetch( `${API_ENDPOINT}/${endpoint}`, opts )
  );
};

/**
 * Retrieves the salt value used to hash the user's password.
 * @param username The name of the user to look up.
 * @returns The salt value (if the user exits).
 */
export const passTheSalt = async ( username: string ) => {
  const response = await buildQuery( 'creds/salt', { username } );
  const { data } = await response.json();

  return data;
};

/**
 * Send the locally generated password hash to the server to authenticate user and request access.
 * @param action Whether to initiate a authenticated session or confirm an existing session.
 * @param hash The locally generated password hash.
 * @param username The email of the user attempting to log in.
 */
export const submitHash = async ( action: TActions, hash: string, username: string ) => {
  const response = await buildQuery( 'guest/auth', {
    action,
    hash,
    username,
  } );

  const res = await response.json();

  console.log( res );
};

export const submitFiles = async ( endpoint: string, body: FormData ) => (
  fetch( `${API_ENDPOINT}/${endpoint}`, {
    body,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    method: 'POST',
  } )
);
