import fetch from 'isomorphic-fetch';

const downloadMovies = (token, conf) => fetch(
  `${conf.BACK_URL}/torrent/add`,
  {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      token,
    }),
  },
);

export default {
  downloadMovies,
};
