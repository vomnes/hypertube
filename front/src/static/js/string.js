const toTitle = (string) => {
  if (string === '') return string;
  return string.charAt(0)
    .toUpperCase() + string.slice(1);
};

const inArray = (category, array) => array.indexOf(category, 0) > -1;

export default {
  toTitle,
  inArray,
};
