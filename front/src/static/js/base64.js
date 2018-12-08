const objectToBase64 = (object) => {
  const ret = Buffer.from(JSON.stringify(object))
    .toString('base64');
  return ret;
};

export default {
  objectToBase64,
};
