const DEFAULT_LIMIT = 10;

/**
 * @returns number
 */
function getCurrentPage() {
  const MIN_PAGE = 1;

  const queryParams = window.location.search;
  if (queryParams === "") {
    return MIN_PAGE;
  }

  const urlParams = new URLSearchParams(window.location.search);
  let offset = urlParams.get("offset");
  if (!offset) {
    return MIN_PAGE;
  } else {
    offset = parseInt(offset);
  }

  let limit = urlParams.get("limit");
  if (!limit) {
    limit = DEFAULT_LIMIT;
  } else {
    limit = parseInt(limit);
  }

  return offset / parseInt(limit) + 1;
}

/**
 * @param {number} pageNumber
 * @returns number
 */
function computeOffset(pageNumber) {
  return DEFAULT_LIMIT * (pageNumber - 1);
}

/**
 * @returns void
 */
function nextPage() {
  const emptyTable = document.getElementById("emptyTable");
  if (emptyTable) return;

  const currPage = getCurrentPage();
  const offset = computeOffset(currPage + 1);
  navigate(offset);
}

/**
 * @returns void
 */
function prevPage() {
  const currPage = getCurrentPage();
  if (currPage <= 1) return;

  const offset = computeOffset(currPage - 1);
  navigate(offset);
}

/**
 * @param {number} offset
 * @returns void
 */
function navigate(offset) {
  const path = window.location.pathname;
  const params = new URLSearchParams(window.location.search);
  const query = params.get("query");

  if (query) {
    window.location.href = `${path}?query=${query}&limit=${DEFAULT_LIMIT}&offset=${offset}`;
  } else {
    window.location.href = `${path}?limit=${DEFAULT_LIMIT}&offset=${offset}`;
  }
}

const pageNumber = document.getElementById("pageNumber");
if (pageNumber) {
  pageNumber.innerText = getCurrentPage();
}
