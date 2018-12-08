const isLoggedIn = () => localStorage.getItem('hypertube_token');

const token = () => localStorage.getItem('hypertube_token');

const storeToken = (value) => {
  localStorage.setItem('hypertube_token', value);
};

const revokeToken = (val) => {
  const http = new XMLHttpRequest();
  http.open('GET', `https://accounts.google.com/o/oauth2/revoke?token=${val}`);
  http.send();
};

const decodeJWT = (jwt) => {
  const data = jwt.split('.');
  if (data.length === 3) {
    try {
      return JSON.parse(atob(data[1]));
    } catch (e) {
      return false;
    }
  }
  return false;
};

const logOut = () => {
  const tokenData = token();
  if (tokenData) {
    const tokenDecoded = decodeJWT(tokenData);
    if (tokenDecoded) {
      if (tokenDecoded.oauth && tokenDecoded.oauth.token && tokenDecoded.oauth.provider === 'gplus') {
        revokeToken(tokenDecoded.oauth.token);
      }
    }
  }
  localStorage.removeItem('hypertube_token');
};

export default {
  isLoggedIn,
  storeToken,
  logOut,
  token,
};
