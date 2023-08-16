/**
 * Converts a string password into a CryptoKey.
 * @param password The password provided by the user.
 * @returns A CryptoKey based on the provided password string.
 */
const generateKey = ( password: string ) => {
  const enc = new TextEncoder();

  return window.crypto.subtle.importKey(
    'raw',
    enc.encode( password ),
    'PBKDF2',
    false,
    ['deriveBits', 'deriveKey'],
  );
};

/**
 * Uses PBKDF2 to hash the user provided password.
 * @param password The password provided by the user.
 * @param salt A random salt value to add to the hash.
 * @returns The hashed password and salt combination.
 */
export const deriveHash = async ( password: string, salt: string ) => {
  const enc = new TextEncoder();

  const keyMaterial = await generateKey( password );

  const dk = await window.crypto.subtle.deriveBits(
    {
      name: 'PBKDF2',
      hash: 'SHA-256',
      iterations: 4096,
      salt: enc.encode( salt ),
    },
    keyMaterial,
    256,
  );

  return btoa( String.fromCharCode( ...new Uint8Array( dk ) ) );
};
