const generateKey = (password: string) => {
  const enc = new TextEncoder();

  return window.crypto.subtle.importKey(
    "raw",
    enc.encode(password),
    "PBKDF2",
    false,
    ["deriveBits", "deriveKey"]
  )
}

export const deriveHash = async (password: string, salt: string) => {
  const enc = new TextEncoder();

  const keyMaterial = await generateKey(password);

  const dk = await window.crypto.subtle.deriveBits(
    {
      name: "PBKDF2",
      hash: "SHA-256",
      iterations: 4096,
      salt: enc.encode(salt),
    },
    keyMaterial,
    256
  )

  const decoded = btoa(String.fromCharCode(...new Uint8Array(dk)));

  return decoded;
}

export const compareHashes = (hash1: string, hash2: string) => {
  if ( hash1 === hash2 ) {
    console.log('MATCH')
  }
}