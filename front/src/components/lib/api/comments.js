import fetch from 'isomorphic-fetch';

const GET = (filmId, token, conf) => fetch(
  `${conf.BACK_URL}/api/v1/comment/${filmId}`,
  {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  },
);

const ADD = (filmId, content, token, conf) => fetch(
  `${conf.BACK_URL}/api/v1/comment/${filmId}`,
  {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`
    },
    body: JSON.stringify({
      content
    })
  }
)

export default {
  GET,
  ADD
};
