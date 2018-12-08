const isScrollXMax = (el) => {
  const ret = (el.scrollLeft === el.scrollLeftMax) ||
    (el.scrollLeft >= (el.scrollWidth - el.clientWidth));
  return ret;
};

const isScrollYMax = container =>
  // == Body ==
  // const ret = el.offsetHeight <= (window.innerHeight + window.scrollY);
  // const ret = el.offsetHeight <= (window.scroll + window.scrollY);
  // return ret;
  container.scrollHeight - container.scrollTop === container.clientHeight
;

export default {
  isScrollXMax,
  isScrollYMax,
};
