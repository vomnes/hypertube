import fetch from 'isomorphic-fetch';

const getProfile = (token, conf) => fetch(
  `${conf.BACK_URL}/api/v1/profiles`,
  {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  },
);

const editData = (token, conf, data) => fetch(
  `${conf.BACK_URL}/api/v1/profiles`,
  {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  },
);

const updatePhoto = (token, conf, data) => fetch(
  `${conf.BACK_URL}/api/v1/profiles/picture`,
  {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  },
);

const updatePassword = (token, conf, data) => fetch(
  `${conf.BACK_URL}/api/v1/profiles/password`,
  {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  },
);

export default {
  getProfile,
  editData,
  updatePhoto,
  updatePassword,
};
