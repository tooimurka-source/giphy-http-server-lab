const form = document.querySelector("#searchForm");
const queryInput = document.querySelector("#queryInput");
const formError = document.querySelector("#formError");
const postSearchButton = document.querySelector("#postSearchButton");
const healthButton = document.querySelector("#healthButton");
const serverStatus = document.querySelector("#serverStatus");
const healthOutput = document.querySelector("#healthOutput");
const responseOutput = document.querySelector("#responseOutput");
const gifPreview = document.querySelector("#gifPreview");
const emptyResult = document.querySelector("#emptyResult");
const lastMethod = document.querySelector("#lastMethod");
const lastQuery = document.querySelector("#lastQuery");

const minQueryLength = 2;
const maxQueryLength = 60;

function getQuery() {
  return queryInput.value.trim();
}

function validateQuery() {
  const query = getQuery();

  if (!query) {
    return "\u0412\u0432\u0435\u0434\u0456\u0442\u044c \u043f\u043e\u0448\u0443\u043a\u043e\u0432\u0438\u0439 \u0437\u0430\u043f\u0438\u0442.";
  }

  if (query.length < minQueryLength) {
    return `\u0417\u0430\u043f\u0438\u0442 \u043c\u0430\u0454 \u043c\u0456\u0441\u0442\u0438\u0442\u0438 \u0449\u043e\u043d\u0430\u0439\u043c\u0435\u043d\u0448\u0435 ${minQueryLength} \u0441\u0438\u043c\u0432\u043e\u043b\u0438.`;
  }

  if (query.length > maxQueryLength) {
    return `\u0417\u0430\u043f\u0438\u0442 \u043c\u0430\u0454 \u043c\u0456\u0441\u0442\u0438\u0442\u0438 \u043d\u0435 \u0431\u0456\u043b\u044c\u0448\u0435 ${maxQueryLength} \u0441\u0438\u043c\u0432\u043e\u043b\u0456\u0432.`;
  }

  return "";
}

function setError(message) {
  formError.textContent = message;
}

function prettyJSON(data) {
  return JSON.stringify(data, null, 2);
}

function setSearching(method, query) {
  gifPreview.removeAttribute("src");
  gifPreview.style.display = "none";
  emptyResult.style.display = "block";
  emptyResult.textContent = "\u0428\u0443\u043a\u0430\u0454\u043c\u043e GIF-\u0430\u043d\u0456\u043c\u0430\u0446\u0456\u044e...";
  lastMethod.textContent = method;
  lastQuery.textContent = query;
  responseOutput.textContent = prettyJSON({ status: "loading", method, query });
}

function setButtonsDisabled(disabled) {
  form.querySelector('button[type="submit"]').disabled = disabled;
  postSearchButton.disabled = disabled;
}

async function parseResponse(response) {
  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "\u0421\u0435\u0440\u0432\u0435\u0440 \u043f\u043e\u0432\u0435\u0440\u043d\u0443\u0432 \u043f\u043e\u043c\u0438\u043b\u043a\u0443.");
  }

  return data;
}

function showGif(data, method) {
  gifPreview.src = `${data.url}#${method}-${Date.now()}`;
  gifPreview.style.display = "block";
  emptyResult.style.display = "none";
  emptyResult.textContent = "GIF-\u0430\u043d\u0456\u043c\u0430\u0446\u0456\u044f \u0437'\u044f\u0432\u0438\u0442\u044c\u0441\u044f \u043f\u0456\u0441\u043b\u044f \u043f\u043e\u0448\u0443\u043a\u0443.";
  lastMethod.textContent = method;
  lastQuery.textContent = data.query;
  responseOutput.textContent = prettyJSON(data);
}

async function searchGif(method) {
  const validationMessage = validateQuery();
  if (validationMessage) {
    setError(validationMessage);
    return;
  }

  setError("");
  const query = getQuery();
  setSearching(method, query);
  setButtonsDisabled(true);

  try {
    const response = method === "GET"
      ? await fetch(`/api/gifs/search?${new URLSearchParams({ q: query })}`)
      : await fetch("/api/gifs/search", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query })
      });

    const data = await parseResponse(response);
    showGif(data, method);
  } finally {
    setButtonsDisabled(false);
  }
}

async function checkHealth() {
  const response = await fetch("/health");
  const data = await parseResponse(response);
  serverStatus.textContent = data.status;
  healthOutput.textContent = prettyJSON(data);
}

function showRequestError(error) {
  setError(error.message);
  responseOutput.textContent = prettyJSON({ error: error.message });
  emptyResult.style.display = "block";
  emptyResult.textContent = "\u041d\u0435 \u0432\u0434\u0430\u043b\u043e\u0441\u044f \u043e\u0442\u0440\u0438\u043c\u0430\u0442\u0438 GIF-\u0430\u043d\u0456\u043c\u0430\u0446\u0456\u044e.";
}

form.addEventListener("submit", (event) => {
  event.preventDefault();
  searchGif("GET").catch(showRequestError);
});

postSearchButton.addEventListener("click", () => {
  searchGif("POST").catch(showRequestError);
});

healthButton.addEventListener("click", () => {
  checkHealth().catch((error) => {
    serverStatus.textContent = "\u043f\u043e\u043c\u0438\u043b\u043a\u0430";
    healthOutput.textContent = prettyJSON({ error: error.message });
  });
});

checkHealth().catch(() => {
  serverStatus.textContent = "\u043d\u0435 \u0432\u0456\u0434\u043f\u043e\u0432\u0456\u0434\u0430\u0454";
});
