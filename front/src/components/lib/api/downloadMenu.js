import fetch from 'isomorphic-fetch';

const torrentList = conf => fetch(
  `${conf.BACK_URL}/torrent/list`,
  {
    method: 'GET',
  },
);

export default {
  torrentList,
};
