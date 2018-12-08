import fetch from 'isomorphic-fetch';

const getMovie = (params, token, conf) => fetch(
  `${conf.BACK_URL}/api/v1/movies/item/${params.filmID}/${params.language}`,
  {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  },
);

const getCategory = (params, token, conf) => fetch(
  `${conf.BACK_URL}/api/v1/movies/category/${params.category}/${params.offset}/${params.numberItems}/${params.language}`,
  {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  },
);

// Parameters to encode in base64
// {
//   search: String,
//   rating: {
//     max: Number,
//     min: Number,
//   },
//   year: {
//     max: Number,
//     min: Number,
//   },
//   status: String,
// }
const searchMovies = (params, token, conf) => fetch(
  `${conf.BACK_URL}/api/v1/movies/search/${params.offset}/${params.numberItems}/${params.language}`,
  {
    method: 'GET',
    headers: {
      'Search-Parameters': params.optionsBase64,
      Authorization: `Bearer ${token}`,
    },
  },
);

const torrents = (params, token, conf) => fetch(
  `${conf.BACK_URL}/api/v1/torrents/${params.imdb}`,
  {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  },
);

const updateViewMovie = (filmId, token, conf) => fetch(
  `${conf.BACK_URL}/api/v1/movies/view/${filmId}`,
  {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  },
);

export default {
  getMovie,
  getCategory,
  searchMovies,
  torrents,
  updateViewMovie
};
