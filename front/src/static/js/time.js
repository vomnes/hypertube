const getHoursFromMins = (mins) => {
  const h = Math.round(mins / 60),
    m = Math.round(mins % 60);
  let minutes = m.toString();
  if (m < 10) {
    minutes = `0${minutes}`;
  }
  return `${h}h${minutes}`;
};

export default {
  getHoursFromMins,
};
