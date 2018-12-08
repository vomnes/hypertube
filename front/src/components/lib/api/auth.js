import fetch from 'isomorphic-fetch';

const login = (params, conf) => fetch(
  `${conf.BACK_URL}/api/v1/accounts/login`,
  {
    method: 'POST',
    body: JSON.stringify({
      username: params.username,
      password: params.password,
    }),
  },
);

const register = (params, conf) => fetch(
  `${conf.BACK_URL}/api/v1/accounts/register`,
  {
    method: 'POST',
    body: JSON.stringify({
      username: params.username,
      email: params.email,
      lastname: params.lastname,
      firstname: params.firstname,
      password: params.password,
      rePassword: params.rePassword,
      picture_base64: params.picture_base64,
    }),
  },
);

export default {
  login,
  register,
};
