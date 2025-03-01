const DEFAULT_LIMIT = 10;

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

function computeOffset(pageNumber) {
  return DEFAULT_LIMIT * (pageNumber - 1);
}

function nextPage() {
  const emptyTable = document.getElementById("emptyTable");
  if (emptyTable) return;

  const currPage = getCurrentPage();
  const offset = computeOffset(currPage + 1);

  const path = window.location.pathname;
  window.location.href = `${path}?limit=${DEFAULT_LIMIT}&offset=${offset}`;
}

function prevPage() {
  const currPage = getCurrentPage();
  if (currPage <= 1) return;

  const offset = computeOffset(currPage - 1);

  const path = window.location.pathname;
  window.location.href = `${path}?limit=${DEFAULT_LIMIT}&offset=${offset}`;
}

const pageNumber = document.getElementById("pageNumber");
if (pageNumber) {
  pageNumber.innerText = getCurrentPage();
}
