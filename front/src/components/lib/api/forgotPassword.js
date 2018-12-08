import fetch from 'isomorphic-fetch';

const sendMail = (params, conf) => fetch(
  `${conf.BACK_URL}/api/v1/mails/forgotpassword`,
  {
    method: 'POST',
    body: JSON.stringify({
      email: params.email,
    }),
  },
);

const change = (params, conf) => fetch(
  `${conf.BACK_URL}/api/v1/accounts/resetpassword`,
  {
    method: 'POST',
    body: JSON.stringify({
      randomToken: params.randomToken,
      password: params.password,
      rePassword: params.rePassword,
    }),
  },
);

export default {
  sendMail,
  change,
};
