import auth from './auth';
import config from './config';
import profile from './profile';
import forgotPassword from './forgotPassword';
import movies from './movies';
import torrent from './downloadMenu';
import dl from './torrentList';
import comments from './comments';

const apiConf = {
  BACK_URL: config.apiDomain,
};

export default {
  // Public
  login: params => auth.login(params, apiConf),
  register: params => auth.register(params, apiConf),
  forgotPasswordSendEmail: params => forgotPassword.sendMail(params, apiConf),
  forgotPasswordChange: params => forgotPassword.change(params, apiConf),
  // Private
  getProfile: token => profile.getProfile(token, apiConf),
  getMovie: (params, token) => movies.getMovie(params, token, apiConf),
  getCategory: (params, token) => movies.getCategory(params, token, apiConf),
  editData: (token, data) => profile.editData(token, apiConf, data),
  searchMovies: (params, token) => movies.searchMovies(params, token, apiConf),
  torrents: (params, token) => movies.torrents(params, token, apiConf),
  updatePhoto: (token, data) => profile.updatePhoto(token, apiConf, data),
  updatePassword: (token, data) => profile.updatePassword(token, apiConf, data),
  torrentList: () => torrent.torrentList(apiConf),
  downloadMovies: data => dl.downloadMovies(data, apiConf),
  getFilmComments: (filmId, token) => comments.GET(filmId, token, apiConf),
  addFilmComments: (filmId, content, token) => comments.ADD(filmId, content, token, apiConf),
  updateViewMovie: (filmId, token) => movies.updateViewMovie(filmId, token, apiConf)
};
