function getCurLanguage() {
  return sessionStorage.getItem('language') ? sessionStorage.getItem('language') : 'en';
}

function setCurLanguage(lang) {
  if (typeof (Storage) !== 'undefined') {
    sessionStorage.setItem('language', lang);
  }
}

module.exports = {
  getCurLanguage,
  setCurLanguage,
};
