export const fetchUserCredentials = async (username: string) => {
  const response = await fetch('', {
    method: 'GET',
    body: JSON.stringify({ username }),
  });

  const { hash, salt } = await response.json();

  return { hash, salt };
};