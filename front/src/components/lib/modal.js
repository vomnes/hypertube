import swal from 'sweetalert';

const errorFormatLanguage = (content, code) => {
  let error = 'accessDenied';
  if (code === 403) {
    error = 'accessDenied';
  }
  // TODO: Fix translation (move into locales.js)
  const translations = {
    usernamePasswordIncorrect: 'Username or password incorrect',
    emptyField: ['Cannot have an empty field', 'At least one field of the body is empty', 'No field inside the body can be empty'],
    invalidUsername: 'Not a valid username',
    invalidFirstname: 'Not a valid firstname',
    invalidLastname: 'Not a valid lastname',
    invalidEmail: ['Not a valid email address', 'Email address is not valid'],
    identicalPassword: 'Both password entered must be identical',
    invalidPassword: 'Not a valid password',
    noPicture: 'Base64 can\'t be empty',
    noEmailDB: 'Email address does not exists in the database',
    invalidData: 'Failed to decode body',
    invalidPicture: ['Corrupted file', 'not accepted, support only png, jpg and jpeg images'],
    emailAlreadyUsed: 'Email address already used',
    currentPasswordInvalid: 'Current password incorrect',
  };
  Object.keys(translations)
    .forEach((key) => {
      if (Array.isArray(translations[key])) {
        Object.keys(translations[key])
          .forEach((element) => {
            if (translations[key][element] === content ||
              content.includes(translations[key][element])) {
              error = key;
            }
          });
      } else if (translations[key] === content) {
        error = key;
      }
    });
  return error;
};

const open = (data) => {
  swal({
    text: data.content,
    icon: data.type,
    dangerMode: data.type === 'error',
  });
};

export default {
  open,
  errorFormatLanguage,
};
